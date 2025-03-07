package main

import (
	"bufio"
	"fmt"
	"image/color"
	"os"
	"strings"

	"github.com/eliukblau/pixterm/pkg/ansimage"
	"github.com/fogleman/gg"
	"golang.org/x/term"
)

func main() {
	p := os.Args[3]
	if err := applyMemefile(p); err != nil {
		panic(err)
	}
}

func applyMemefile(p string) error {
	mf, err := parseMemefile(p)
	if err != nil {
		panic(err)
	}
	ctx := gg.NewContext(0, 0)
	for _, c := range mf {
		if err := c.Apply(ctx); err != nil {
			return err
		}

	}

	w, h, err := term.GetSize(0)
	if err != nil {
		// Default to standard terminal size
		w, h = 80, 24
	}

	ai, err := ansimage.NewScaledFromImage(ctx.Image(), h, w, color.Transparent, ansimage.ScaleModeFit, ansimage.NoDithering)
	if err != nil {
		return err
	}
	ai.Draw()
	return ctx.SavePNG("output.png")
}

func parseMemefile(p string) (Memefile, error) {
	fmt.Println("Reading file", p)
	fmt.Println("---")
	res := []command{}
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	s.Split(bufio.ScanLines)

	for s.Scan() {
		line := s.Text()
		if line == "" {
			continue
		}
		before, after, found := strings.Cut(line, " ")
		if !found {
			continue
		}
		fmt.Println(before, after)
		var add command
		switch before {
		case "FROM":
			add = &fromCmd{source: after}
		case "TOP":
			add = &topCmd{text: after}
		case "BOTTOM":
			add = &bottomCmd{text: after}
		}
		res = append(res, add)
	}
	fmt.Println("---")
	return res, nil
}

type Memefile []command

type command interface {
	Apply(ctx *gg.Context) error
}

type topCmd struct {
	text string
}

func (c *topCmd) Apply(ctx *gg.Context) error {
	TopBanner(ctx, c.text)
	return nil
}

type bottomCmd struct {
	text string
}

func (c *bottomCmd) Apply(ctx *gg.Context) error {
	BottomBanner(ctx, c.text)
	return nil
}

type fromCmd struct {
	source string
}

func (c *fromCmd) Apply(ctx *gg.Context) error {
	// this replaces the ctx
	img, err := gg.LoadImage(c.source)
	if err != nil {
		return err
	}
	*ctx = *NewContext(img)
	return nil
}
