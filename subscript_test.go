package subscript

import (
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/testutil"
)

type TestCase struct {
	desc string
	md   string
	html string
}

func TestGoldmarkOnly(t *testing.T) {
	// These tests are to show how Goldmark handles strikethrough by default,
	// without our extension enabled.
	mdTest := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.DefinitionList,
			extension.Footnote,
			// NewSubscript(),
		),
	)

	testCases := []TestCase{
		{
			desc: "Goldmark only: single-tilde strikethrough",
			md: `H~2~O`,
			html: `<p>H<del>2</del>O</p>`,
		},
		{
			desc: "Goldmark only: single-tilde and dbl-tilde strikethrough",
			md: `H~2~O with a ~~del~~ strikethrough`,
			html: `<p>H<del>2</del>O with a <del>del</del> strikethrough</p>`,
		},
		{
			desc: "Goldmark only: glucose formula",
			md: `C~6~H~12~O~6~`,
			html: `<p>C<del>6</del>H<del>12</del>O<del>6</del></p>`,
		},
		{
			desc: "Goldmark only: glucose formula with dbl-tilde strikethrough",
			md: `C~6~H~12~O~6~ is ~~not~~ critical for life`,
			html: `<p>C<del>6</del>H<del>12</del>O<del>6</del> is <del>not</del> critical for life</p>`,
		},
		{
			desc: "Goldmark only: glucose formula with single-tilde strikethrough with spaces",
			md: `C~6~H~12~O~6~ ~is not~ is critical for life`,
			html: `<p>C<del>6</del>H<del>12</del>O<del>6</del> <del>is not</del> is critical for life</p>`,
		},
		{
			desc: "Goldmark only: glucose formula with spaces inner, leading and trailing spaces",
			md: `C~6 0~H~12 a1 ~O~ 6 ~ ~is not~ is critical for life`,
			html: `<p>C<del>6 0</del>H~12 a1 <del>O</del> 6 ~ <del>is not</del> is critical for life</p>`,
		},
		{
			desc: "Goldmark only: glucose formula with emphasis (strong)",
			md: `**C~6~H~12~O~6~**`,
			html: `<p><strong>C<del>6</del>H<del>12</del>O<del>6</del></strong></p>`,
		},
		{
			desc: "Goldmark only: no trailing spaces inside strikethrough",
			md: `H~2 ~O`,
			html: `<p>H~2 ~O</p>`,
		},
		{
			desc: "Goldmark only: no leading spaces inside strikethrough",
			md: `H~  2~O`,
			html: `<p>H~  2~O</p>`,
		},
		{
			desc: "Goldmark only: no leading spaces inside strikethrough",
			md: `H~2 abc~O`,
			html: `<p>H<del>2 abc</del>O</p>`,
		},
		// {
		// 	desc: "",
		// 	md: ``,
		// 	html: ``,
		// },
		// {
		// 	desc: "",
		// 	md: ``,
		// 	html: ``,
		// },
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			testutil.DoTestCase(mdTest, testutil.MarkdownTestCase{
				Description: tc.desc,
				Markdown:    tc.md,
				Expected:    tc.html,
			}, t)
		})
	}

}

func TestSubscriptCore(t *testing.T) {
	// Because extension.GFM enables strikethrough by default, we need to include it here.
	// If we don't, we won't be doing a reliable test to make sure subscript and
	// strikethrough work together properly.
	mdTest := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.DefinitionList,
			extension.Footnote,
			NewSubscript(),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)

	testCases := []TestCase{
		{
			desc: "Subscript: basic test",
			md: `H~2~O`,
			html: `<p>H<sub>2</sub>O</p>`,
		},
		{
			desc: "Subscript: basic test with dbl-tilde strikethrough",
			md: `H~2~O with a ~~del~~ strikethrough`,
			html: `<p>H<sub>2</sub>O with a <del>del</del> strikethrough</p>`,
		},
		{
			desc: "Subscript: glucose formula",
			md: `C~6~H~12~O~6~`,
			html: `<p>C<sub>6</sub>H<sub>12</sub>O<sub>6</sub></p>`,
		},
		{
			desc: "Subscript: glucose formula with dbl-tilde strikethrough",
			md: `C~6~H~12~O~6~ is ~~not~~ critical for life`,
			html: `<p>C<sub>6</sub>H<sub>12</sub>O<sub>6</sub> is <del>not</del> critical for life</p>`,
		},
		{
			desc: "Subscript: glucose formula with single-tilde strikethrough",
			md: `C~6~H~12~O~6~ is ~not~ critical for life`,
			html: `<p>C<sub>6</sub>H<sub>12</sub>O<sub>6</sub> is <del>not</del> critical for life</p>`,
		},
		{
			desc: "Subscript: glucose formula with single-tilde strikethrough with spaces",
			md: `C~6~H~12~O~6~ ~is not~ is critical for life`,
			html: `<p>C<sub>6</sub>H<sub>12</sub>O<sub>6</sub> <del>is not</del> is critical for life</p>`,
		},
		{
			desc: "Subscript: glucose formula with single-tilde strikethrough with spaces",
			md: `C~6~H~12~O~6~ ~~is not~~ is critical for life`,
			html: `<p>C<sub>6</sub>H<sub>12</sub>O<sub>6</sub> <del>is not</del> is critical for life</p>`,
		},
		{
			desc: "Subscript: glucose formula with dbl-tilde strikethrough",
			md: `C~6~H~12~O~6~ is ~~not~~ critical for life`,
			html: `<p>C<sub>6</sub>H<sub>12</sub>O<sub>6</sub> is <del>not</del> critical for life</p>`,
		},
		{
			desc: "Subscript: glucose formula with emphasis (strong)",
			md: `**C~6~H~12~O~6~** is **critical** for life`,
			html: `<p><strong>C<sub>6</sub>H<sub>12</sub>O<sub>6</sub></strong> is <strong>critical</strong> for life</p>`,
		},
		{
			desc: "Subscript: no trailing spaces inside subscript OR strikethrough",
			md: `H~2 ~O`,
			html: `<p>H~2 ~O</p>`,
		},
		{
			desc: "Subscript: no leading spaces inside subscript OR strikethrough",
			md: `H~ 2~O`,
			html: `<p>H~ 2~O</p>`,
		},
		{
			desc: "Subscript: no spaces inside subscript (treat as strikethrough)",
			md: `H~2 abc~O`,
			html: `<p>H<del>2 abc</del>O</p>`,
		},
		{
			desc: "Subscript: subscript cannot start at beginning of line",
			md: `~2~O`,
			html: `<p><del>2</del>O</p>`,
		},
		{
			desc: "Subscript: subscript must have non-whitespace before it",
			md: `H ~2~ O`,
			html: `<p>H <del>2</del> O</p>`,
		},
		{
			// NOTE: This is the correct output!
			// The subscript parser runs first, so any valid subscripts WILL be processed as subscripts.
			// Because neither strikethrough nor subscript allow leading or trailing spaces inside their delimiters,
			// any tildes with a space before or after it is just treated as a normal character.
			// Since THIS test ends with ' ~', the final tilde is just treated as a normal character.
			desc: "Subscript: glucose formula with spaces inside subscript",
			md: `C~6 ~H~ 12~O~ 6 ~`,
			html: `<p>C~6 <del>H</del> 12<sub>O</sub> 6 ~</p>`,
		},
		{
			// NOTE: This is the correct output!
			// The subscript parser runs first, so any valid subscripts WILL be processed as subscripts.
			// The only difference between this and the last test case is that this one ends with '6~' instead of '6 ~'.
			// Since the space is no longer between the final '6' and the final '~', and the first tilde after the 'C'
			// at the beginning of the line is followed by a non-space (the '6'), Goldmark will process everything,
			// including the previously rendered subscripts and strikethroughs as another strikethrough for the
			// entire line EXCEPT the 'C' at the beginning of the line.
			desc: "Subscript: glucose formula with spaces inside subscript",
			md: `C~6 ~H~ 12~O~ 6~`,
			html: `<p>C<del>6 <del>H</del> 12<sub>O</sub> 6</del></p>`,
		},
		{
			// NOTE: This is the correct output!
			// Subscripts cannot have spaces inside them,
			// so this gets treated similar to the last test case, except the outer double-tilde strikethrough
			// gets processed and strikes through everything inside it, even the valid subscript.
			desc: "Subscript: glucose formula with spaces inside subscript",
			md: `~~C~6 ~H~ 12~O~ 6~~`,
			html: `<p><del>C~6 <del>H</del> 12<sub>O</sub> 6</del></p>`,
		},
		{
			desc: "Subscript: other characters allowed inside subscript",
			md: `Foo~1,2~ + Bar~(test)~ - Baz~[abc]~ * Quux~{xyz}~ < Zzz~<tag>~`,
			html: `<p>Foo<sub>1,2</sub> + Bar<sub>(test)</sub> - Baz<sub>[abc]</sub> * Quux<sub>{xyz}</sub> &lt; Zzz<sub>&lt;tag&gt;</sub></p>`,
		},
		{
			desc: "Subscript: Drake equation",
			md: `- N = R~&#x1f7af;~ x f~p~ x n~e~ x f~l~ x f~i~ x f~c~ x L`,
			// md: `- N = R~&#x2731;~ x f~p~ x n~e~ x f~l~ x f~i~ x f~c~ x L`,
			html: ``,
		},
		// {
		// 	desc: "",
		// 	md: ``,
		// 	html: ``,
		// },
		// {
		// 	desc: "",
		// 	md: ``,
		// 	html: ``,
		// },
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			testutil.DoTestCase(mdTest, testutil.MarkdownTestCase{
				Description: tc.desc,
				Markdown:    tc.md,
				Expected:    tc.html,
			}, t)
		})
	}

}