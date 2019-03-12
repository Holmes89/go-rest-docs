package gorestdoc

import "strings"

type MarkDownBuilder struct {
	builder strings.Builder
}

func (b *MarkDownBuilder) Build() string {
	return b.builder.String()
}

func (b *MarkDownBuilder) H1(title string) *MarkDownBuilder {
	b.builder.WriteString("# ")
	b.builder.WriteString(title)
	b.builder.WriteString("\n")
	return b
}

func (b *MarkDownBuilder) H2(title string) *MarkDownBuilder {
	b.builder.WriteString("## ")
	b.builder.WriteString(title)
	b.builder.WriteString("\n")
	return b
}

func (b *MarkDownBuilder) H3(title string) *MarkDownBuilder {
	b.builder.WriteString("### ")
	b.builder.WriteString(title)
	b.builder.WriteString("\n")
	return b
}

func (b *MarkDownBuilder) H4(title string) *MarkDownBuilder {
	b.builder.WriteString("#### ")
	b.builder.WriteString(title)
	b.builder.WriteString("\n")
	return b
}

func (b *MarkDownBuilder) Body(body string) *MarkDownBuilder {
	b.builder.WriteString(body)
	b.builder.WriteString("\n\n")
	return b
}

func (b *MarkDownBuilder) Code(code string) *MarkDownBuilder {
	b.builder.WriteString("```\n")
	b.builder.WriteString(code)
	b.builder.WriteString("\n```\n\n")
	return b
}
