package ast

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

func ifaceStr(val interface{}) string {
	str := ""
	for _, seg := range val.([]interface{}) {
		str = str + string(seg.([]byte))
	}

	return str
}

func ifaceSequences(val interface{}) []Sequence {
	sentences := []Sequence{}
	for _, node := range val.([]interface{}) {
		sentences = append(sentences, node.(Sequence))
	}

	return sentences
}

func ifaceNodes(val interface{}) []Node {
	nodes := []Node{}
	for _, node := range val.([]interface{}) {
		nodes = append(nodes, node.(Node))
	}

	return nodes
}

var g = &grammar{
	rules: []*rule{
		{
			name: "Booklit",
			pos:  position{line: 34, col: 1, offset: 605},
			expr: &actionExpr{
				pos: position{line: 34, col: 12, offset: 616},
				run: (*parser).callonBooklit1,
				expr: &seqExpr{
					pos: position{line: 34, col: 12, offset: 616},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 34, col: 12, offset: 616},
							label: "node",
							expr: &ruleRefExpr{
								pos:  position{line: 34, col: 17, offset: 621},
								name: "Paragraphs",
							},
						},
						&notExpr{
							pos: position{line: 34, col: 28, offset: 632},
							expr: &anyMatcher{
								line: 34, col: 29, offset: 633,
							},
						},
					},
				},
			},
		},
		{
			name: "Paragraphs",
			pos:  position{line: 38, col: 1, offset: 659},
			expr: &actionExpr{
				pos: position{line: 38, col: 15, offset: 673},
				run: (*parser).callonParagraphs1,
				expr: &seqExpr{
					pos: position{line: 38, col: 15, offset: 673},
					exprs: []interface{}{
						&zeroOrMoreExpr{
							pos: position{line: 38, col: 15, offset: 673},
							expr: &litMatcher{
								pos:        position{line: 38, col: 15, offset: 673},
								val:        "\n",
								ignoreCase: false,
							},
						},
						&labeledExpr{
							pos:   position{line: 38, col: 21, offset: 679},
							label: "paragraphs",
							expr: &oneOrMoreExpr{
								pos: position{line: 38, col: 32, offset: 690},
								expr: &actionExpr{
									pos: position{line: 38, col: 33, offset: 691},
									run: (*parser).callonParagraphs7,
									expr: &seqExpr{
										pos: position{line: 38, col: 33, offset: 691},
										exprs: []interface{}{
											&labeledExpr{
												pos:   position{line: 38, col: 33, offset: 691},
												label: "p",
												expr: &ruleRefExpr{
													pos:  position{line: 38, col: 35, offset: 693},
													name: "Paragraph",
												},
											},
											&zeroOrMoreExpr{
												pos: position{line: 38, col: 45, offset: 703},
												expr: &litMatcher{
													pos:        position{line: 38, col: 45, offset: 703},
													val:        "\n",
													ignoreCase: false,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Paragraph",
			pos:  position{line: 42, col: 1, offset: 781},
			expr: &actionExpr{
				pos: position{line: 42, col: 14, offset: 794},
				run: (*parser).callonParagraph1,
				expr: &labeledExpr{
					pos:   position{line: 42, col: 14, offset: 794},
					label: "sentences",
					expr: &oneOrMoreExpr{
						pos: position{line: 42, col: 24, offset: 804},
						expr: &actionExpr{
							pos: position{line: 42, col: 25, offset: 805},
							run: (*parser).callonParagraph4,
							expr: &seqExpr{
								pos: position{line: 42, col: 25, offset: 805},
								exprs: []interface{}{
									&labeledExpr{
										pos:   position{line: 42, col: 25, offset: 805},
										label: "s",
										expr: &ruleRefExpr{
											pos:  position{line: 42, col: 27, offset: 807},
											name: "Sentence",
										},
									},
									&litMatcher{
										pos:        position{line: 42, col: 36, offset: 816},
										val:        "\n",
										ignoreCase: false,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Sentence",
			pos:  position{line: 46, col: 1, offset: 897},
			expr: &actionExpr{
				pos: position{line: 46, col: 13, offset: 909},
				run: (*parser).callonSentence1,
				expr: &seqExpr{
					pos: position{line: 46, col: 13, offset: 909},
					exprs: []interface{}{
						&zeroOrMoreExpr{
							pos: position{line: 46, col: 13, offset: 909},
							expr: &charClassMatcher{
								pos:        position{line: 46, col: 13, offset: 909},
								val:        "[ \\t]",
								chars:      []rune{' ', '\t'},
								ignoreCase: false,
								inverted:   false,
							},
						},
						&labeledExpr{
							pos:   position{line: 46, col: 20, offset: 916},
							label: "words",
							expr: &oneOrMoreExpr{
								pos: position{line: 46, col: 26, offset: 922},
								expr: &ruleRefExpr{
									pos:  position{line: 46, col: 27, offset: 923},
									name: "Word",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Word",
			pos:  position{line: 50, col: 1, offset: 977},
			expr: &actionExpr{
				pos: position{line: 50, col: 9, offset: 985},
				run: (*parser).callonWord1,
				expr: &labeledExpr{
					pos:   position{line: 50, col: 9, offset: 985},
					label: "val",
					expr: &choiceExpr{
						pos: position{line: 50, col: 14, offset: 990},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 50, col: 14, offset: 990},
								name: "String",
							},
							&ruleRefExpr{
								pos:  position{line: 50, col: 23, offset: 999},
								name: "Invoke",
							},
						},
					},
				},
			},
		},
		{
			name: "SplitSentence",
			pos:  position{line: 54, col: 1, offset: 1030},
			expr: &actionExpr{
				pos: position{line: 54, col: 18, offset: 1047},
				run: (*parser).callonSplitSentence1,
				expr: &seqExpr{
					pos: position{line: 54, col: 18, offset: 1047},
					exprs: []interface{}{
						&zeroOrMoreExpr{
							pos: position{line: 54, col: 18, offset: 1047},
							expr: &charClassMatcher{
								pos:        position{line: 54, col: 18, offset: 1047},
								val:        "[ \\t]",
								chars:      []rune{' ', '\t'},
								ignoreCase: false,
								inverted:   false,
							},
						},
						&labeledExpr{
							pos:   position{line: 54, col: 25, offset: 1054},
							label: "firstWord",
							expr: &ruleRefExpr{
								pos:  position{line: 54, col: 35, offset: 1064},
								name: "Word",
							},
						},
						&labeledExpr{
							pos:   position{line: 54, col: 40, offset: 1069},
							label: "words",
							expr: &zeroOrMoreExpr{
								pos: position{line: 54, col: 46, offset: 1075},
								expr: &choiceExpr{
									pos: position{line: 54, col: 47, offset: 1076},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 54, col: 47, offset: 1076},
											name: "Word",
										},
										&ruleRefExpr{
											pos:  position{line: 54, col: 54, offset: 1083},
											name: "Split",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Split",
			pos:  position{line: 59, col: 1, offset: 1214},
			expr: &actionExpr{
				pos: position{line: 59, col: 10, offset: 1223},
				run: (*parser).callonSplit1,
				expr: &litMatcher{
					pos:        position{line: 59, col: 10, offset: 1223},
					val:        "\n",
					ignoreCase: false,
				},
			},
		},
		{
			name: "String",
			pos:  position{line: 61, col: 1, offset: 1257},
			expr: &choiceExpr{
				pos: position{line: 61, col: 11, offset: 1267},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 61, col: 11, offset: 1267},
						run: (*parser).callonString2,
						expr: &labeledExpr{
							pos:   position{line: 61, col: 11, offset: 1267},
							label: "str",
							expr: &oneOrMoreExpr{
								pos: position{line: 61, col: 15, offset: 1271},
								expr: &charClassMatcher{
									pos:        position{line: 61, col: 15, offset: 1271},
									val:        "[^\\\\{}\\n]",
									chars:      []rune{'\\', '{', '}', '\n'},
									ignoreCase: false,
									inverted:   true,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 61, col: 59, offset: 1315},
						run: (*parser).callonString6,
						expr: &seqExpr{
							pos: position{line: 61, col: 59, offset: 1315},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 61, col: 59, offset: 1315},
									val:        "\\",
									ignoreCase: false,
								},
								&charClassMatcher{
									pos:        position{line: 61, col: 64, offset: 1320},
									val:        "[\\\\{}]",
									chars:      []rune{'\\', '{', '}'},
									ignoreCase: false,
									inverted:   false,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "VerbatimString",
			pos:  position{line: 63, col: 1, offset: 1363},
			expr: &choiceExpr{
				pos: position{line: 63, col: 19, offset: 1381},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 63, col: 19, offset: 1381},
						run: (*parser).callonVerbatimString2,
						expr: &labeledExpr{
							pos:   position{line: 63, col: 19, offset: 1381},
							label: "str",
							expr: &oneOrMoreExpr{
								pos: position{line: 63, col: 23, offset: 1385},
								expr: &charClassMatcher{
									pos:        position{line: 63, col: 23, offset: 1385},
									val:        "[^\\n}]",
									chars:      []rune{'\n', '}'},
									ignoreCase: false,
									inverted:   true,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 65, col: 5, offset: 1428},
						run: (*parser).callonVerbatimString6,
						expr: &seqExpr{
							pos: position{line: 65, col: 5, offset: 1428},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 65, col: 5, offset: 1428},
									val:        "}",
									ignoreCase: false,
								},
								&notExpr{
									pos: position{line: 65, col: 9, offset: 1432},
									expr: &litMatcher{
										pos:        position{line: 65, col: 10, offset: 1433},
										val:        "}}",
										ignoreCase: false,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Indent",
			pos:  position{line: 69, col: 1, offset: 1472},
			expr: &actionExpr{
				pos: position{line: 69, col: 11, offset: 1482},
				run: (*parser).callonIndent1,
				expr: &zeroOrMoreExpr{
					pos: position{line: 69, col: 11, offset: 1482},
					expr: &charClassMatcher{
						pos:        position{line: 69, col: 11, offset: 1482},
						val:        "[ \\t]",
						chars:      []rune{' ', '\t'},
						ignoreCase: false,
						inverted:   false,
					},
				},
			},
		},
		{
			name: "PreformattedSentence",
			pos:  position{line: 86, col: 1, offset: 1747},
			expr: &actionExpr{
				pos: position{line: 86, col: 25, offset: 1771},
				run: (*parser).callonPreformattedSentence1,
				expr: &seqExpr{
					pos: position{line: 86, col: 25, offset: 1771},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 86, col: 25, offset: 1771},
							label: "indent",
							expr: &ruleRefExpr{
								pos:  position{line: 86, col: 32, offset: 1778},
								name: "Indent",
							},
						},
						&labeledExpr{
							pos:   position{line: 86, col: 39, offset: 1785},
							label: "words",
							expr: &zeroOrMoreExpr{
								pos: position{line: 86, col: 45, offset: 1791},
								expr: &choiceExpr{
									pos: position{line: 86, col: 46, offset: 1792},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 86, col: 46, offset: 1792},
											name: "String",
										},
										&ruleRefExpr{
											pos:  position{line: 86, col: 55, offset: 1801},
											name: "Invoke",
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 86, col: 64, offset: 1810},
							val:        "\n",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Preformatted",
			pos:  position{line: 92, col: 1, offset: 1935},
			expr: &actionExpr{
				pos: position{line: 92, col: 17, offset: 1951},
				run: (*parser).callonPreformatted1,
				expr: &seqExpr{
					pos: position{line: 92, col: 17, offset: 1951},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 92, col: 17, offset: 1951},
							val:        "\n",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 92, col: 22, offset: 1956},
							label: "sentences",
							expr: &zeroOrMoreExpr{
								pos: position{line: 92, col: 32, offset: 1966},
								expr: &ruleRefExpr{
									pos:  position{line: 92, col: 32, offset: 1966},
									name: "PreformattedSentence",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "VerbatimSentence",
			pos:  position{line: 97, col: 1, offset: 2086},
			expr: &actionExpr{
				pos: position{line: 97, col: 21, offset: 2106},
				run: (*parser).callonVerbatimSentence1,
				expr: &seqExpr{
					pos: position{line: 97, col: 21, offset: 2106},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 97, col: 21, offset: 2106},
							label: "indent",
							expr: &ruleRefExpr{
								pos:  position{line: 97, col: 28, offset: 2113},
								name: "Indent",
							},
						},
						&labeledExpr{
							pos:   position{line: 97, col: 35, offset: 2120},
							label: "words",
							expr: &zeroOrMoreExpr{
								pos: position{line: 97, col: 41, offset: 2126},
								expr: &ruleRefExpr{
									pos:  position{line: 97, col: 41, offset: 2126},
									name: "VerbatimString",
								},
							},
						},
						&litMatcher{
							pos:        position{line: 97, col: 57, offset: 2142},
							val:        "\n",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Verbatim",
			pos:  position{line: 103, col: 1, offset: 2267},
			expr: &actionExpr{
				pos: position{line: 103, col: 13, offset: 2279},
				run: (*parser).callonVerbatim1,
				expr: &seqExpr{
					pos: position{line: 103, col: 13, offset: 2279},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 103, col: 13, offset: 2279},
							val:        "\n",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 103, col: 18, offset: 2284},
							label: "sentences",
							expr: &zeroOrMoreExpr{
								pos: position{line: 103, col: 28, offset: 2294},
								expr: &ruleRefExpr{
									pos:  position{line: 103, col: 28, offset: 2294},
									name: "VerbatimSentence",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Invoke",
			pos:  position{line: 108, col: 1, offset: 2410},
			expr: &actionExpr{
				pos: position{line: 108, col: 11, offset: 2420},
				run: (*parser).callonInvoke1,
				expr: &seqExpr{
					pos: position{line: 108, col: 11, offset: 2420},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 108, col: 11, offset: 2420},
							val:        "\\",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 108, col: 16, offset: 2425},
							label: "name",
							expr: &oneOrMoreExpr{
								pos: position{line: 108, col: 22, offset: 2431},
								expr: &charClassMatcher{
									pos:        position{line: 108, col: 22, offset: 2431},
									val:        "[a-z-]",
									chars:      []rune{'-'},
									ranges:     []rune{'a', 'z'},
									ignoreCase: false,
									inverted:   false,
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 108, col: 31, offset: 2440},
							label: "args",
							expr: &zeroOrMoreExpr{
								pos: position{line: 108, col: 37, offset: 2446},
								expr: &ruleRefExpr{
									pos:  position{line: 108, col: 37, offset: 2446},
									name: "Argument",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "VerbatimArg",
			pos:  position{line: 115, col: 1, offset: 2551},
			expr: &actionExpr{
				pos: position{line: 115, col: 16, offset: 2566},
				run: (*parser).callonVerbatimArg1,
				expr: &seqExpr{
					pos: position{line: 115, col: 16, offset: 2566},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 115, col: 16, offset: 2566},
							val:        "{{{",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 115, col: 22, offset: 2572},
							label: "node",
							expr: &ruleRefExpr{
								pos:  position{line: 115, col: 27, offset: 2577},
								name: "Verbatim",
							},
						},
						&zeroOrMoreExpr{
							pos: position{line: 115, col: 36, offset: 2586},
							expr: &charClassMatcher{
								pos:        position{line: 115, col: 36, offset: 2586},
								val:        "[ \\t]",
								chars:      []rune{' ', '\t'},
								ignoreCase: false,
								inverted:   false,
							},
						},
						&litMatcher{
							pos:        position{line: 115, col: 43, offset: 2593},
							val:        "}}}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "PreformattedArg",
			pos:  position{line: 119, col: 1, offset: 2623},
			expr: &actionExpr{
				pos: position{line: 119, col: 20, offset: 2642},
				run: (*parser).callonPreformattedArg1,
				expr: &seqExpr{
					pos: position{line: 119, col: 20, offset: 2642},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 119, col: 20, offset: 2642},
							val:        "{{",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 119, col: 25, offset: 2647},
							label: "node",
							expr: &ruleRefExpr{
								pos:  position{line: 119, col: 30, offset: 2652},
								name: "Preformatted",
							},
						},
						&zeroOrMoreExpr{
							pos: position{line: 119, col: 43, offset: 2665},
							expr: &charClassMatcher{
								pos:        position{line: 119, col: 43, offset: 2665},
								val:        "[ \\t]",
								chars:      []rune{' ', '\t'},
								ignoreCase: false,
								inverted:   false,
							},
						},
						&litMatcher{
							pos:        position{line: 119, col: 50, offset: 2672},
							val:        "}}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Arg",
			pos:  position{line: 123, col: 1, offset: 2701},
			expr: &actionExpr{
				pos: position{line: 123, col: 8, offset: 2708},
				run: (*parser).callonArg1,
				expr: &seqExpr{
					pos: position{line: 123, col: 8, offset: 2708},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 123, col: 8, offset: 2708},
							val:        "{",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 123, col: 12, offset: 2712},
							label: "node",
							expr: &choiceExpr{
								pos: position{line: 123, col: 18, offset: 2718},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 123, col: 18, offset: 2718},
										name: "SplitSentence",
									},
									&ruleRefExpr{
										pos:  position{line: 123, col: 34, offset: 2734},
										name: "Paragraphs",
									},
								},
							},
						},
						&zeroOrMoreExpr{
							pos: position{line: 123, col: 46, offset: 2746},
							expr: &charClassMatcher{
								pos:        position{line: 123, col: 46, offset: 2746},
								val:        "[ \\t]",
								chars:      []rune{' ', '\t'},
								ignoreCase: false,
								inverted:   false,
							},
						},
						&litMatcher{
							pos:        position{line: 123, col: 53, offset: 2753},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Argument",
			pos:  position{line: 127, col: 1, offset: 2781},
			expr: &choiceExpr{
				pos: position{line: 127, col: 13, offset: 2793},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 127, col: 13, offset: 2793},
						name: "VerbatimArg",
					},
					&ruleRefExpr{
						pos:  position{line: 127, col: 27, offset: 2807},
						name: "PreformattedArg",
					},
					&ruleRefExpr{
						pos:  position{line: 127, col: 45, offset: 2825},
						name: "Arg",
					},
				},
			},
		},
	},
}

func (c *current) onBooklit1(node interface{}) (interface{}, error) {
	return node, nil
}

func (p *parser) callonBooklit1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBooklit1(stack["node"])
}

func (c *current) onParagraphs7(p interface{}) (interface{}, error) {
	return p, nil
}

func (p *parser) callonParagraphs7() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onParagraphs7(stack["p"])
}

func (c *current) onParagraphs1(paragraphs interface{}) (interface{}, error) {
	return Sequence(ifaceNodes(paragraphs)), nil
}

func (p *parser) callonParagraphs1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onParagraphs1(stack["paragraphs"])
}

func (c *current) onParagraph4(s interface{}) (interface{}, error) {
	return s, nil
}

func (p *parser) callonParagraph4() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onParagraph4(stack["s"])
}

func (c *current) onParagraph1(sentences interface{}) (interface{}, error) {
	return Paragraph(ifaceSequences(sentences)), nil
}

func (p *parser) callonParagraph1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onParagraph1(stack["sentences"])
}

func (c *current) onSentence1(words interface{}) (interface{}, error) {
	return Sequence(ifaceNodes(words)), nil
}

func (p *parser) callonSentence1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSentence1(stack["words"])
}

func (c *current) onWord1(val interface{}) (interface{}, error) {
	return val, nil
}

func (p *parser) callonWord1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onWord1(stack["val"])
}

func (c *current) onSplitSentence1(firstWord, words interface{}) (interface{}, error) {
	allWords := append([]interface{}{firstWord}, words.([]interface{})...)
	return Sequence(ifaceNodes(allWords)), nil
}

func (p *parser) callonSplitSentence1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSplitSentence1(stack["firstWord"], stack["words"])
}

func (c *current) onSplit1() (interface{}, error) {
	return String(" "), nil
}

func (p *parser) callonSplit1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSplit1()
}

func (c *current) onString2(str interface{}) (interface{}, error) {
	return String(c.text), nil
}

func (p *parser) callonString2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onString2(stack["str"])
}

func (c *current) onString6() (interface{}, error) {
	return String(c.text[1:]), nil
}

func (p *parser) callonString6() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onString6()
}

func (c *current) onVerbatimString2(str interface{}) (interface{}, error) {
	return String(c.text), nil
}

func (p *parser) callonVerbatimString2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onVerbatimString2(stack["str"])
}

func (c *current) onVerbatimString6() (interface{}, error) {
	return String(c.text), nil
}

func (p *parser) callonVerbatimString6() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onVerbatimString6()
}

func (c *current) onIndent1() (interface{}, error) {
	skip := len(c.text)

	i, found := c.globalStore["indent-skip"]
	if found {
		skip = i.(int)
	} else {
		c.globalStore["indent-skip"] = skip
	}

	if skip <= len(c.text) {
		return string(c.text[skip:]), nil
	} else {
		return "", nil
	}
}

func (p *parser) callonIndent1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIndent1()
}

func (c *current) onPreformattedSentence1(indent, words interface{}) (interface{}, error) {
	line := []Node{String(indent.(string))}
	line = append(line, ifaceNodes(words)...)
	return Sequence(line), nil
}

func (p *parser) callonPreformattedSentence1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPreformattedSentence1(stack["indent"], stack["words"])
}

func (c *current) onPreformatted1(sentences interface{}) (interface{}, error) {
	delete(c.globalStore, "indentation")
	return Preformatted(ifaceSequences(sentences)), nil
}

func (p *parser) callonPreformatted1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPreformatted1(stack["sentences"])
}

func (c *current) onVerbatimSentence1(indent, words interface{}) (interface{}, error) {
	line := []Node{String(indent.(string))}
	line = append(line, ifaceNodes(words)...)
	return Sequence(line), nil
}

func (p *parser) callonVerbatimSentence1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onVerbatimSentence1(stack["indent"], stack["words"])
}

func (c *current) onVerbatim1(sentences interface{}) (interface{}, error) {
	delete(c.globalStore, "indentation")
	return Preformatted(ifaceSequences(sentences)), nil
}

func (p *parser) callonVerbatim1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onVerbatim1(stack["sentences"])
}

func (c *current) onInvoke1(name, args interface{}) (interface{}, error) {
	return Invoke{
		Function:  ifaceStr(name),
		Arguments: ifaceNodes(args),
	}, nil
}

func (p *parser) callonInvoke1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onInvoke1(stack["name"], stack["args"])
}

func (c *current) onVerbatimArg1(node interface{}) (interface{}, error) {
	return node, nil
}

func (p *parser) callonVerbatimArg1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onVerbatimArg1(stack["node"])
}

func (c *current) onPreformattedArg1(node interface{}) (interface{}, error) {
	return node, nil
}

func (p *parser) callonPreformattedArg1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPreformattedArg1(stack["node"])
}

func (c *current) onArg1(node interface{}) (interface{}, error) {
	return node, nil
}

func (p *parser) callonArg1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onArg1(stack["node"])
}

var (
	// errNoRule is returned when the grammar to parse has no rule.
	errNoRule = errors.New("grammar has no rule")

	// errInvalidEncoding is returned when the source is not properly
	// utf8-encoded.
	errInvalidEncoding = errors.New("invalid encoding")
)

// Option is a function that can set an option on the parser. It returns
// the previous setting as an Option.
type Option func(*parser) Option

// Debug creates an Option to set the debug flag to b. When set to true,
// debugging information is printed to stdout while parsing.
//
// The default is false.
func Debug(b bool) Option {
	return func(p *parser) Option {
		old := p.debug
		p.debug = b
		return Debug(old)
	}
}

// Memoize creates an Option to set the memoize flag to b. When set to true,
// the parser will cache all results so each expression is evaluated only
// once. This guarantees linear parsing time even for pathological cases,
// at the expense of more memory and slower times for typical cases.
//
// The default is false.
func Memoize(b bool) Option {
	return func(p *parser) Option {
		old := p.memoize
		p.memoize = b
		return Memoize(old)
	}
}

// Recover creates an Option to set the recover flag to b. When set to
// true, this causes the parser to recover from panics and convert it
// to an error. Setting it to false can be useful while debugging to
// access the full stack trace.
//
// The default is true.
func Recover(b bool) Option {
	return func(p *parser) Option {
		old := p.recover
		p.recover = b
		return Recover(old)
	}
}

// GlobalStore creates an Option to set a key to a certain value in
// the globalStore.
func GlobalStore(key string, value interface{}) Option {
	return func(p *parser) Option {
		old := p.cur.globalStore[key]
		p.cur.globalStore[key] = value
		return GlobalStore(key, old)
	}
}

// ParseFile parses the file identified by filename.
func ParseFile(filename string, opts ...Option) (i interface{}, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = f.Close()
	}()
	return ParseReader(filename, f, opts...)
}

// ParseReader parses the data from r using filename as information in the
// error messages.
func ParseReader(filename string, r io.Reader, opts ...Option) (interface{}, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return Parse(filename, b, opts...)
}

// Parse parses the data from b using filename as information in the
// error messages.
func Parse(filename string, b []byte, opts ...Option) (interface{}, error) {
	return newParser(filename, b, opts...).parse(g)
}

// position records a position in the text.
type position struct {
	line, col, offset int
}

func (p position) String() string {
	return fmt.Sprintf("%d:%d [%d]", p.line, p.col, p.offset)
}

// savepoint stores all state required to go back to this point in the
// parser.
type savepoint struct {
	position
	rn rune
	w  int
}

type current struct {
	pos  position // start position of the match
	text []byte   // raw text of the match

	// the globalStore allows the parser to store arbitrary values
	globalStore map[string]interface{}
}

// the AST types...

type grammar struct {
	pos   position
	rules []*rule
}

type rule struct {
	pos         position
	name        string
	displayName string
	expr        interface{}
}

type choiceExpr struct {
	pos          position
	alternatives []interface{}
}

type actionExpr struct {
	pos  position
	expr interface{}
	run  func(*parser) (interface{}, error)
}

type seqExpr struct {
	pos   position
	exprs []interface{}
}

type labeledExpr struct {
	pos   position
	label string
	expr  interface{}
}

type expr struct {
	pos  position
	expr interface{}
}

type andExpr expr
type notExpr expr
type zeroOrOneExpr expr
type zeroOrMoreExpr expr
type oneOrMoreExpr expr

type ruleRefExpr struct {
	pos  position
	name string
}

type andCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type notCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type litMatcher struct {
	pos        position
	val        string
	ignoreCase bool
}

type charClassMatcher struct {
	pos        position
	val        string
	chars      []rune
	ranges     []rune
	classes    []*unicode.RangeTable
	ignoreCase bool
	inverted   bool
}

type anyMatcher position

// errList cumulates the errors found by the parser.
type errList []error

func (e *errList) add(err error) {
	*e = append(*e, err)
}

func (e errList) err() error {
	if len(e) == 0 {
		return nil
	}
	e.dedupe()
	return e
}

func (e *errList) dedupe() {
	var cleaned []error
	set := make(map[string]bool)
	for _, err := range *e {
		if msg := err.Error(); !set[msg] {
			set[msg] = true
			cleaned = append(cleaned, err)
		}
	}
	*e = cleaned
}

func (e errList) Error() string {
	switch len(e) {
	case 0:
		return ""
	case 1:
		return e[0].Error()
	default:
		var buf bytes.Buffer

		for i, err := range e {
			if i > 0 {
				buf.WriteRune('\n')
			}
			buf.WriteString(err.Error())
		}
		return buf.String()
	}
}

// parserError wraps an error with a prefix indicating the rule in which
// the error occurred. The original error is stored in the Inner field.
type parserError struct {
	Inner    error
	pos      position
	prefix   string
	expected []string
}

// Error returns the error message.
func (p *parserError) Error() string {
	return p.prefix + ": " + p.Inner.Error()
}

// newParser creates a parser with the specified input source and options.
func newParser(filename string, b []byte, opts ...Option) *parser {
	p := &parser{
		filename: filename,
		errs:     new(errList),
		data:     b,
		pt:       savepoint{position: position{line: 1}},
		recover:  true,
		cur: current{
			globalStore: make(map[string]interface{}),
		},
		maxFailPos:      position{col: 1, line: 1},
		maxFailExpected: make(map[string]struct{}),
	}
	p.setOptions(opts)
	return p
}

// setOptions applies the options to the parser.
func (p *parser) setOptions(opts []Option) {
	for _, opt := range opts {
		opt(p)
	}
}

type resultTuple struct {
	v   interface{}
	b   bool
	end savepoint
}

type parser struct {
	filename string
	pt       savepoint
	cur      current

	data []byte
	errs *errList

	depth   int
	recover bool
	debug   bool

	memoize bool
	// memoization table for the packrat algorithm:
	// map[offset in source] map[expression or rule] {value, match}
	memo map[int]map[interface{}]resultTuple

	// rules table, maps the rule identifier to the rule node
	rules map[string]*rule
	// variables stack, map of label to value
	vstack []map[string]interface{}
	// rule stack, allows identification of the current rule in errors
	rstack []*rule

	// stats
	exprCnt int

	// parse fail
	maxFailPos            position
	maxFailExpected       map[string]struct{}
	maxFailInvertExpected bool
}

// push a variable set on the vstack.
func (p *parser) pushV() {
	if cap(p.vstack) == len(p.vstack) {
		// create new empty slot in the stack
		p.vstack = append(p.vstack, nil)
	} else {
		// slice to 1 more
		p.vstack = p.vstack[:len(p.vstack)+1]
	}

	// get the last args set
	m := p.vstack[len(p.vstack)-1]
	if m != nil && len(m) == 0 {
		// empty map, all good
		return
	}

	m = make(map[string]interface{})
	p.vstack[len(p.vstack)-1] = m
}

// pop a variable set from the vstack.
func (p *parser) popV() {
	// if the map is not empty, clear it
	m := p.vstack[len(p.vstack)-1]
	if len(m) > 0 {
		// GC that map
		p.vstack[len(p.vstack)-1] = nil
	}
	p.vstack = p.vstack[:len(p.vstack)-1]
}

func (p *parser) print(prefix, s string) string {
	if !p.debug {
		return s
	}

	fmt.Printf("%s %d:%d:%d: %s [%#U]\n",
		prefix, p.pt.line, p.pt.col, p.pt.offset, s, p.pt.rn)
	return s
}

func (p *parser) in(s string) string {
	p.depth++
	return p.print(strings.Repeat(" ", p.depth)+">", s)
}

func (p *parser) out(s string) string {
	p.depth--
	return p.print(strings.Repeat(" ", p.depth)+"<", s)
}

func (p *parser) addErr(err error) {
	p.addErrAt(err, p.pt.position, []string{})
}

func (p *parser) addErrAt(err error, pos position, expected []string) {
	var buf bytes.Buffer
	if p.filename != "" {
		buf.WriteString(p.filename)
	}
	if buf.Len() > 0 {
		buf.WriteString(":")
	}
	buf.WriteString(fmt.Sprintf("%d:%d (%d)", pos.line, pos.col, pos.offset))
	if len(p.rstack) > 0 {
		if buf.Len() > 0 {
			buf.WriteString(": ")
		}
		rule := p.rstack[len(p.rstack)-1]
		if rule.displayName != "" {
			buf.WriteString("rule " + rule.displayName)
		} else {
			buf.WriteString("rule " + rule.name)
		}
	}
	pe := &parserError{Inner: err, pos: pos, prefix: buf.String(), expected: expected}
	p.errs.add(pe)
}

func (p *parser) failAt(fail bool, pos position, want string) {
	// process fail if parsing fails and not inverted or parsing succeeds and invert is set
	if fail == p.maxFailInvertExpected {
		if pos.offset < p.maxFailPos.offset {
			return
		}

		if pos.offset > p.maxFailPos.offset {
			p.maxFailPos = pos
			p.maxFailExpected = make(map[string]struct{})
		}

		if p.maxFailInvertExpected {
			want = "!" + want
		}
		p.maxFailExpected[want] = struct{}{}
	}
}

// read advances the parser to the next rune.
func (p *parser) read() {
	p.pt.offset += p.pt.w
	rn, n := utf8.DecodeRune(p.data[p.pt.offset:])
	p.pt.rn = rn
	p.pt.w = n
	p.pt.col++
	if rn == '\n' {
		p.pt.line++
		p.pt.col = 0
	}

	if rn == utf8.RuneError {
		if n == 1 {
			p.addErr(errInvalidEncoding)
		}
	}
}

// restore parser position to the savepoint pt.
func (p *parser) restore(pt savepoint) {
	if p.debug {
		defer p.out(p.in("restore"))
	}
	if pt.offset == p.pt.offset {
		return
	}
	p.pt = pt
}

// get the slice of bytes from the savepoint start to the current position.
func (p *parser) sliceFrom(start savepoint) []byte {
	return p.data[start.position.offset:p.pt.position.offset]
}

func (p *parser) getMemoized(node interface{}) (resultTuple, bool) {
	if len(p.memo) == 0 {
		return resultTuple{}, false
	}
	m := p.memo[p.pt.offset]
	if len(m) == 0 {
		return resultTuple{}, false
	}
	res, ok := m[node]
	return res, ok
}

func (p *parser) setMemoized(pt savepoint, node interface{}, tuple resultTuple) {
	if p.memo == nil {
		p.memo = make(map[int]map[interface{}]resultTuple)
	}
	m := p.memo[pt.offset]
	if m == nil {
		m = make(map[interface{}]resultTuple)
		p.memo[pt.offset] = m
	}
	m[node] = tuple
}

func (p *parser) buildRulesTable(g *grammar) {
	p.rules = make(map[string]*rule, len(g.rules))
	for _, r := range g.rules {
		p.rules[r.name] = r
	}
}

func (p *parser) parse(g *grammar) (val interface{}, err error) {
	if len(g.rules) == 0 {
		p.addErr(errNoRule)
		return nil, p.errs.err()
	}

	// TODO : not super critical but this could be generated
	p.buildRulesTable(g)

	if p.recover {
		// panic can be used in action code to stop parsing immediately
		// and return the panic as an error.
		defer func() {
			if e := recover(); e != nil {
				if p.debug {
					defer p.out(p.in("panic handler"))
				}
				val = nil
				switch e := e.(type) {
				case error:
					p.addErr(e)
				default:
					p.addErr(fmt.Errorf("%v", e))
				}
				err = p.errs.err()
			}
		}()
	}

	// start rule is rule [0]
	p.read() // advance to first rune
	val, ok := p.parseRule(g.rules[0])
	if !ok {
		if len(*p.errs) == 0 {
			// If parsing fails, but no errors have been recorded, the expected values
			// for the farthest parser position are returned as error.
			expected := make([]string, 0, len(p.maxFailExpected))
			eof := false
			if _, ok := p.maxFailExpected["!."]; ok {
				delete(p.maxFailExpected, "!.")
				eof = true
			}
			for k := range p.maxFailExpected {
				expected = append(expected, k)
			}
			sort.Strings(expected)
			if eof {
				expected = append(expected, "EOF")
			}
			p.addErrAt(errors.New("no match found, expected: "+listJoin(expected, ", ", "or")), p.maxFailPos, expected)
		}
		return nil, p.errs.err()
	}
	return val, p.errs.err()
}

func listJoin(list []string, sep string, lastSep string) string {
	switch len(list) {
	case 0:
		return ""
	case 1:
		return list[0]
	default:
		return fmt.Sprintf("%s %s %s", strings.Join(list[:len(list)-1], sep), lastSep, list[len(list)-1])
	}
}

func (p *parser) parseRule(rule *rule) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRule " + rule.name))
	}

	if p.memoize {
		res, ok := p.getMemoized(rule)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
	}

	start := p.pt
	p.rstack = append(p.rstack, rule)
	p.pushV()
	val, ok := p.parseExpr(rule.expr)
	p.popV()
	p.rstack = p.rstack[:len(p.rstack)-1]
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}

	if p.memoize {
		p.setMemoized(start, rule, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseExpr(expr interface{}) (interface{}, bool) {
	var pt savepoint

	if p.memoize {
		res, ok := p.getMemoized(expr)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
		pt = p.pt
	}

	p.exprCnt++
	var val interface{}
	var ok bool
	switch expr := expr.(type) {
	case *actionExpr:
		val, ok = p.parseActionExpr(expr)
	case *andCodeExpr:
		val, ok = p.parseAndCodeExpr(expr)
	case *andExpr:
		val, ok = p.parseAndExpr(expr)
	case *anyMatcher:
		val, ok = p.parseAnyMatcher(expr)
	case *charClassMatcher:
		val, ok = p.parseCharClassMatcher(expr)
	case *choiceExpr:
		val, ok = p.parseChoiceExpr(expr)
	case *labeledExpr:
		val, ok = p.parseLabeledExpr(expr)
	case *litMatcher:
		val, ok = p.parseLitMatcher(expr)
	case *notCodeExpr:
		val, ok = p.parseNotCodeExpr(expr)
	case *notExpr:
		val, ok = p.parseNotExpr(expr)
	case *oneOrMoreExpr:
		val, ok = p.parseOneOrMoreExpr(expr)
	case *ruleRefExpr:
		val, ok = p.parseRuleRefExpr(expr)
	case *seqExpr:
		val, ok = p.parseSeqExpr(expr)
	case *zeroOrMoreExpr:
		val, ok = p.parseZeroOrMoreExpr(expr)
	case *zeroOrOneExpr:
		val, ok = p.parseZeroOrOneExpr(expr)
	default:
		panic(fmt.Sprintf("unknown expression type %T", expr))
	}
	if p.memoize {
		p.setMemoized(pt, expr, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseActionExpr(act *actionExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseActionExpr"))
	}

	start := p.pt
	val, ok := p.parseExpr(act.expr)
	if ok {
		p.cur.pos = start.position
		p.cur.text = p.sliceFrom(start)
		actVal, err := act.run(p)
		if err != nil {
			p.addErrAt(err, start.position, []string{})
		}
		val = actVal
	}
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}
	return val, ok
}

func (p *parser) parseAndCodeExpr(and *andCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndCodeExpr"))
	}

	ok, err := and.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, ok
}

func (p *parser) parseAndExpr(and *andExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(and.expr)
	p.popV()
	p.restore(pt)
	return nil, ok
}

func (p *parser) parseAnyMatcher(any *anyMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAnyMatcher"))
	}

	if p.pt.rn != utf8.RuneError {
		start := p.pt
		p.read()
		p.failAt(true, start.position, ".")
		return p.sliceFrom(start), true
	}
	p.failAt(false, p.pt.position, ".")
	return nil, false
}

func (p *parser) parseCharClassMatcher(chr *charClassMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseCharClassMatcher"))
	}

	cur := p.pt.rn
	start := p.pt
	// can't match EOF
	if cur == utf8.RuneError {
		p.failAt(false, start.position, chr.val)
		return nil, false
	}
	if chr.ignoreCase {
		cur = unicode.ToLower(cur)
	}

	// try to match in the list of available chars
	for _, rn := range chr.chars {
		if rn == cur {
			if chr.inverted {
				p.failAt(false, start.position, chr.val)
				return nil, false
			}
			p.read()
			p.failAt(true, start.position, chr.val)
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of ranges
	for i := 0; i < len(chr.ranges); i += 2 {
		if cur >= chr.ranges[i] && cur <= chr.ranges[i+1] {
			if chr.inverted {
				p.failAt(false, start.position, chr.val)
				return nil, false
			}
			p.read()
			p.failAt(true, start.position, chr.val)
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of Unicode classes
	for _, cl := range chr.classes {
		if unicode.Is(cl, cur) {
			if chr.inverted {
				p.failAt(false, start.position, chr.val)
				return nil, false
			}
			p.read()
			p.failAt(true, start.position, chr.val)
			return p.sliceFrom(start), true
		}
	}

	if chr.inverted {
		p.read()
		p.failAt(true, start.position, chr.val)
		return p.sliceFrom(start), true
	}
	p.failAt(false, start.position, chr.val)
	return nil, false
}

func (p *parser) parseChoiceExpr(ch *choiceExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseChoiceExpr"))
	}

	for _, alt := range ch.alternatives {
		p.pushV()
		val, ok := p.parseExpr(alt)
		p.popV()
		if ok {
			return val, ok
		}
	}
	return nil, false
}

func (p *parser) parseLabeledExpr(lab *labeledExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLabeledExpr"))
	}

	p.pushV()
	val, ok := p.parseExpr(lab.expr)
	p.popV()
	if ok && lab.label != "" {
		m := p.vstack[len(p.vstack)-1]
		m[lab.label] = val
	}
	return val, ok
}

func (p *parser) parseLitMatcher(lit *litMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLitMatcher"))
	}

	ignoreCase := ""
	if lit.ignoreCase {
		ignoreCase = "i"
	}
	val := fmt.Sprintf("%q%s", lit.val, ignoreCase)
	start := p.pt
	for _, want := range lit.val {
		cur := p.pt.rn
		if lit.ignoreCase {
			cur = unicode.ToLower(cur)
		}
		if cur != want {
			p.failAt(false, start.position, val)
			p.restore(start)
			return nil, false
		}
		p.read()
	}
	p.failAt(true, start.position, val)
	return p.sliceFrom(start), true
}

func (p *parser) parseNotCodeExpr(not *notCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotCodeExpr"))
	}

	ok, err := not.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, !ok
}

func (p *parser) parseNotExpr(not *notExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotExpr"))
	}

	pt := p.pt
	p.pushV()
	p.maxFailInvertExpected = !p.maxFailInvertExpected
	_, ok := p.parseExpr(not.expr)
	p.maxFailInvertExpected = !p.maxFailInvertExpected
	p.popV()
	p.restore(pt)
	return nil, !ok
}

func (p *parser) parseOneOrMoreExpr(expr *oneOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseOneOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			if len(vals) == 0 {
				// did not match once, no match
				return nil, false
			}
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseRuleRefExpr(ref *ruleRefExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRuleRefExpr " + ref.name))
	}

	if ref.name == "" {
		panic(fmt.Sprintf("%s: invalid rule: missing name", ref.pos))
	}

	rule := p.rules[ref.name]
	if rule == nil {
		p.addErr(fmt.Errorf("undefined rule: %s", ref.name))
		return nil, false
	}
	return p.parseRule(rule)
}

func (p *parser) parseSeqExpr(seq *seqExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseSeqExpr"))
	}

	var vals []interface{}

	pt := p.pt
	for _, expr := range seq.exprs {
		val, ok := p.parseExpr(expr)
		if !ok {
			p.restore(pt)
			return nil, false
		}
		vals = append(vals, val)
	}
	return vals, true
}

func (p *parser) parseZeroOrMoreExpr(expr *zeroOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseZeroOrOneExpr(expr *zeroOrOneExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrOneExpr"))
	}

	p.pushV()
	val, _ := p.parseExpr(expr.expr)
	p.popV()
	// whether it matched or not, consider it a match
	return val, true
}