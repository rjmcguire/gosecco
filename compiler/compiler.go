package compiler

import (
	"github.com/twtiger/go-seccomp/constants"
	"github.com/twtiger/go-seccomp/tree"
	"golang.org/x/sys/unix"
)

func newCompiler() *compiler {
	return &compiler{
		currentlyLoaded: -1,
		positiveLabels:  make(map[string][]uint),
		negativeLabels:  make(map[string][]uint),
	}
}

func Compile(policy tree.Policy) ([]unix.SockFilter, error) {
	c := newCompiler()
	c.compile(policy.Rules)
	return c.result, nil
}

type compiler struct {
	result          []unix.SockFilter
	currentlyLoaded int
	positiveLabels  map[string][]uint
	negativeLabels  map[string][]uint
}

func (c *compiler) compile(rules []tree.Rule) {
	for _, r := range rules {
		c.compileRule(r)
	}
	c.positiveAction("")
	c.negativeAction("")
}

func (c *compiler) labelHere(label string) {
	c.fixupJumpPoints(label, uint(len(c.result)))
}

func (c *compiler) compileRule(r tree.Rule) {
	c.labelHere("next")
	c.checkCorrectSyscall(r.Name)
}

const syscallNameIndex = 0

const LOAD = BPF_LD | BPF_W | BPF_ABS
const JEQ_K = BPF_JMP | BPF_JEQ | BPF_K
const RET_K = BPF_RET | BPF_K

func (c *compiler) op(code uint16, k uint32) uint {
	ix := uint(len(c.result))
	c.result = append(c.result, unix.SockFilter{
		Code: code,
		Jt:   0,
		Jf:   0,
		K:    k,
	})
	return ix
}

func (c *compiler) loadAt(pos uint32) {
	if c.currentlyLoaded != int(pos) {
		c.op(LOAD, pos)
		c.currentlyLoaded = int(pos)
	}
}

func (c *compiler) loadCurrentSyscall() {
	c.loadAt(syscallNameIndex)
}

func (c *compiler) positiveJumpTo(index uint, label string) {
	if label != "" {
		c.positiveLabels[label] = append(c.positiveLabels[label], index)
	}
}

func (c *compiler) negativeJumpTo(index uint, label string) {
	if label != "" {
		c.negativeLabels[label] = append(c.negativeLabels[label], index)
	}
}

func (c *compiler) jumpIfEqualTo(val uint32, jt, jf string) {
	num := c.op(JEQ_K, val)
	c.positiveJumpTo(num, jt)
	c.negativeJumpTo(num, jf)
}

func (c *compiler) checkCorrectSyscall(name string) {
	sys, ok := constants.GetSyscall(name)
	if !ok {
		panic("This shouldn't happen - analyzer should have caught it before compiler tries to compile it")
	}

	c.loadCurrentSyscall()
	c.jumpIfEqualTo(sys, "positiveResult", "next")
}

func (c *compiler) fixupJumpPoints(label string, ix uint) {
	for _, origin := range c.positiveLabels[label] {
		// TODO: check that these jumps aren't to large - in that case we need to insert a JUMP_K instruction
		c.result[origin].Jt = uint8(ix-origin) - 1
	}
	delete(c.positiveLabels, label)

	for _, origin := range c.negativeLabels[label] {
		// TODO: check that these jumps aren't to large - in that case we need to insert a JUMP_K instruction
		c.result[origin].Jf = uint8(ix-origin) - 1
	}
	delete(c.negativeLabels, label)
}

func (c *compiler) positiveAction(name string) {
	c.labelHere("positiveResult")
	c.op(RET_K, SECCOMP_RET_ALLOW)
}

func (c *compiler) negativeAction(name string) {
	c.labelHere("next")
	c.op(RET_K, SECCOMP_RET_KILL)
}
