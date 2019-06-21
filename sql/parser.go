// Copyright 2019 The SQLFlow Authors. All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by goyacc -p sql -o parser.go sql.y. DO NOT EDIT.

//line sql.y:2
package sql

import __yyfmt__ "fmt"

//line sql.y:2

import (
	"fmt"
	"strings"
	"sync"
)

/* expr defines an expression as a Lisp list.  If len(val)>0,
   it is an atomic expression, in particular, NUMBER, IDENT,
   or STRING, defined by typ and val; otherwise, it is a
   Lisp S-expression. */
type expr struct {
	typ  int
	val  string
	sexp exprlist
}

type exprlist []*expr

/* construct an atomic expr */
func atomic(typ int, val string) *expr {
	return &expr{
		typ: typ,
		val: val,
	}
}

/* construct a funcall expr */
func funcall(name string, oprd exprlist) *expr {
	return &expr{
		sexp: append(exprlist{atomic(IDENT, name)}, oprd...),
	}
}

/* construct a unary expr */
func unary(typ int, op string, od1 *expr) *expr {
	return &expr{
		sexp: append(exprlist{atomic(typ, op)}, od1),
	}
}

/* construct a binary expr */
func binary(typ int, od1 *expr, op string, od2 *expr) *expr {
	return &expr{
		sexp: append(exprlist{atomic(typ, op)}, od1, od2),
	}
}

/* construct a variadic expr */
func variadic(typ int, op string, ods exprlist) *expr {
	return &expr{
		sexp: append(exprlist{atomic(typ, op)}, ods...),
	}
}

type extendedSelect struct {
	extended bool
	train    bool
	standardSelect
	trainClause
	predictClause
}

type standardSelect struct {
	fields []string
	tables []string
	where  *expr
	limit  string
}

type trainClause struct {
	estimator string
	attrs     attrs
	columns   columnClause
	label     string
	save      string
}

/* If no FOR in the COLUMN, the key is "" */
type columnClause map[string]exprlist

type attrs map[string]*expr

type predictClause struct {
	model string
	into  string
}

var parseResult *extendedSelect

func attrsUnion(as1, as2 attrs) attrs {
	for k, v := range as2 {
		if _, ok := as1[k]; ok {
			log.Panicf("attr %q already specified", as2)
		}
		as1[k] = v
	}
	return as1
}

//line sql.y:104
type sqlSymType struct {
	yys  int
	val  string /* NUMBER, IDENT, STRING, and keywords */
	flds []string
	tbls []string
	expr *expr
	expl exprlist
	atrs attrs
	eslt extendedSelect
	slct standardSelect
	tran trainClause
	colc columnClause
	labc string
	infr predictClause
}

const SELECT = 57346
const FROM = 57347
const WHERE = 57348
const LIMIT = 57349
const TRAIN = 57350
const PREDICT = 57351
const WITH = 57352
const COLUMN = 57353
const LABEL = 57354
const USING = 57355
const INTO = 57356
const FOR = 57357
const IDENT = 57358
const NUMBER = 57359
const STRING = 57360
const AND = 57361
const OR = 57362
const GE = 57363
const LE = 57364
const NOT = 57365
const POWER = 57366
const UMINUS = 57367

var sqlToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"SELECT",
	"FROM",
	"WHERE",
	"LIMIT",
	"TRAIN",
	"PREDICT",
	"WITH",
	"COLUMN",
	"LABEL",
	"USING",
	"INTO",
	"FOR",
	"IDENT",
	"NUMBER",
	"STRING",
	"AND",
	"OR",
	"'>'",
	"'<'",
	"'='",
	"GE",
	"LE",
	"'+'",
	"'-'",
	"'*'",
	"'/'",
	"'%'",
	"NOT",
	"POWER",
	"UMINUS",
	"';'",
	"','",
	"'('",
	"')'",
	"'['",
	"']'",
}
var sqlStatenames = [...]string{}

const sqlEofCode = 1
const sqlErrCode = 2
const sqlInitialStackSize = 16

//line sql.y:271

/* Like Lisp's builtin function cdr. */
func (e *expr) cdr() (r []string) {
	for i := 1; i < len(e.sexp); i++ {
		r = append(r, e.sexp[i].String())
	}
	return r
}

func (e *expr) String() string {
	if e.typ == 0 { /* a compound expression */
		switch e.sexp[0].typ {
		case '+', '*', '/', '%', '=', '<', '>', LE, GE, AND, OR:
			if len(e.sexp) != 3 {
				log.Panicf("Expecting binary expression, got %.10q", e.sexp)
			}
			return fmt.Sprintf("%s %s %s", e.sexp[1], e.sexp[0].val, e.sexp[2])
		case '-':
			switch len(e.sexp) {
			case 2:
				return fmt.Sprintf(" -%s", e.sexp[1])
			case 3:
				return fmt.Sprintf("%s - %s", e.sexp[1], e.sexp[2])
			default:
				log.Panicf("Expecting either unary or binary -, got %.10q", e.sexp)
			}
		case '(':
			if len(e.sexp) != 2 {
				log.Panicf("Expecting ( ) as unary operator, got %.10q", e.sexp)
			}
			return fmt.Sprintf("(%s)", e.sexp[1])
		case '[':
			return "[" + strings.Join(e.cdr(), ", ") + "]"
		case NOT:
			return fmt.Sprintf("NOT %s", e.sexp[1])
		case IDENT: /* function call */
			return e.sexp[0].val + "(" + strings.Join(e.cdr(), ", ") + ")"
		}
	} else {
		return fmt.Sprintf("%s", e.val)
	}

	log.Panicf("Cannot print an unknown expression")
	return ""
}

func (s standardSelect) String() string {
	r := "SELECT "
	if len(s.fields) == 0 {
		r += "*"
	} else {
		r += strings.Join(s.fields, ", ")
	}
	r += "\nFROM " + strings.Join(s.tables, ", ")
	if s.where != nil {
		r += fmt.Sprintf("\nWHERE %s", s.where)
	}
	if len(s.limit) > 0 {
		r += fmt.Sprintf("\nLIMIT %s", s.limit)
	}
	return r + ";"
}

// sqlReentrantParser makes sqlParser, generated by goyacc and using a
// global variable parseResult to return the result, reentrant.
type sqlSyncParser struct {
	pr sqlParser
}

func newParser() *sqlSyncParser {
	return &sqlSyncParser{sqlNewParser()}
}

var mu sync.Mutex

func (p *sqlSyncParser) Parse(s string) (r *extendedSelect, e error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			e, ok = r.(error)
			if !ok {
				e = fmt.Errorf("%v", r)
			}
		}
	}()

	mu.Lock()
	defer mu.Unlock()

	p.pr.Parse(newLexer(s))
	return parseResult, nil
}

func (e *expr) label() string {
	return e.val[1 : len(e.val)-1]
}

//line yacctab:1
var sqlExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const sqlPrivate = 57344

const sqlLast = 152

var sqlAct = [...]int{

	26, 52, 89, 88, 51, 75, 72, 72, 101, 78,
	73, 20, 46, 44, 45, 41, 40, 39, 43, 42,
	34, 35, 36, 37, 38, 33, 32, 47, 99, 48,
	49, 71, 16, 82, 15, 83, 57, 58, 59, 60,
	61, 62, 63, 64, 65, 66, 67, 68, 22, 21,
	23, 70, 7, 9, 8, 10, 11, 81, 19, 28,
	98, 54, 104, 27, 36, 37, 38, 96, 25, 97,
	29, 50, 91, 102, 79, 100, 76, 22, 21, 23,
	99, 4, 77, 92, 90, 93, 92, 87, 28, 95,
	56, 14, 27, 22, 21, 23, 55, 25, 69, 29,
	92, 31, 103, 13, 28, 30, 18, 94, 27, 85,
	86, 53, 3, 25, 74, 29, 44, 45, 41, 40,
	39, 43, 42, 34, 35, 36, 37, 38, 41, 40,
	39, 43, 42, 34, 35, 36, 37, 38, 34, 35,
	36, 37, 38, 24, 17, 12, 6, 84, 80, 5,
	2, 1,
}
var sqlPact = [...]int{

	108, -1000, 47, 75, -1000, 0, -2, 90, 41, 77,
	89, 85, -9, -1000, -1000, -1000, -1000, -10, -1000, -1000,
	97, -1000, -24, -1000, -1000, 77, -1000, 77, 77, 32,
	101, 48, 80, 74, 77, 77, 77, 77, 77, 77,
	77, 77, 77, 77, 77, 77, 61, -6, -1000, -1000,
	-1000, -29, 97, 60, 66, -1000, -1000, 36, 36, -1000,
	-1000, -1000, 112, 112, 112, 112, 112, 107, 107, -1000,
	-28, -1000, 77, -1000, 22, -1000, 12, -1000, -1000, 97,
	98, 60, 56, 77, 93, 56, 51, -1000, 45, -1000,
	-1000, -24, -1000, 97, 59, -7, -1000, -1000, 57, 56,
	-1000, 46, -1000, -1000, -1000,
}
var sqlPgo = [...]int{

	0, 151, 150, 149, 148, 147, 146, 145, 144, 1,
	0, 2, 4, 143, 3, 5, 114,
}
var sqlR1 = [...]int{

	0, 1, 1, 1, 2, 2, 2, 2, 3, 6,
	4, 4, 4, 7, 7, 7, 11, 11, 11, 14,
	14, 5, 5, 8, 8, 15, 16, 16, 10, 10,
	12, 12, 13, 13, 9, 9, 9, 9, 9, 9,
	9, 9, 9, 9, 9, 9, 9, 9, 9, 9,
	9, 9, 9, 9,
}
var sqlR2 = [...]int{

	0, 2, 3, 3, 2, 3, 3, 3, 8, 4,
	2, 4, 5, 1, 1, 3, 1, 1, 1, 1,
	3, 2, 2, 1, 3, 3, 1, 3, 3, 4,
	1, 3, 2, 3, 1, 1, 1, 1, 3, 1,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 2, 2,
}
var sqlChk = [...]int{

	-1000, -1, -2, 4, 34, -3, -6, 5, 7, 6,
	8, 9, -7, 28, 16, 34, 34, -8, 16, 17,
	-9, 17, 16, 18, -13, 36, -10, 31, 27, 38,
	16, 16, 35, 35, 26, 27, 28, 29, 30, 23,
	22, 21, 25, 24, 19, 20, 36, -9, -9, -9,
	39, -12, -9, 10, 13, 16, 16, -9, -9, -9,
	-9, -9, -9, -9, -9, -9, -9, -9, -9, 37,
	-12, 37, 35, 39, -16, -15, 16, 16, 37, -9,
	-4, 35, 11, 23, -5, 11, 12, -15, -14, -11,
	28, 16, -10, -9, 14, -14, 16, 18, 15, 35,
	16, 15, 16, -11, 16,
}
var sqlDef = [...]int{

	0, -2, 0, 0, 1, 0, 0, 0, 0, 0,
	0, 0, 4, 13, 14, 2, 3, 5, 23, 6,
	7, 34, 35, 36, 37, 0, 39, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 52, 53,
	32, 0, 30, 0, 0, 15, 24, 40, 41, 42,
	43, 44, 45, 46, 47, 48, 49, 50, 51, 28,
	0, 38, 0, 33, 0, 26, 0, 9, 29, 31,
	0, 0, 0, 0, 0, 0, 0, 27, 10, 19,
	16, 17, 18, 25, 0, 0, 21, 22, 0, 0,
	8, 0, 11, 20, 12,
}
var sqlTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 30, 3, 3,
	36, 37, 28, 26, 35, 27, 3, 29, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 34,
	22, 23, 21, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 38, 3, 39,
}
var sqlTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 24,
	25, 31, 32, 33,
}
var sqlTok3 = [...]int{
	0,
}

var sqlErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	sqlDebug        = 0
	sqlErrorVerbose = false
)

type sqlLexer interface {
	Lex(lval *sqlSymType) int
	Error(s string)
}

type sqlParser interface {
	Parse(sqlLexer) int
	Lookahead() int
}

type sqlParserImpl struct {
	lval  sqlSymType
	stack [sqlInitialStackSize]sqlSymType
	char  int
}

func (p *sqlParserImpl) Lookahead() int {
	return p.char
}

func sqlNewParser() sqlParser {
	return &sqlParserImpl{}
}

const sqlFlag = -1000

func sqlTokname(c int) string {
	if c >= 1 && c-1 < len(sqlToknames) {
		if sqlToknames[c-1] != "" {
			return sqlToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func sqlStatname(s int) string {
	if s >= 0 && s < len(sqlStatenames) {
		if sqlStatenames[s] != "" {
			return sqlStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func sqlErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !sqlErrorVerbose {
		return "syntax error"
	}

	for _, e := range sqlErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + sqlTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := sqlPact[state]
	for tok := TOKSTART; tok-1 < len(sqlToknames); tok++ {
		if n := base + tok; n >= 0 && n < sqlLast && sqlChk[sqlAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if sqlDef[state] == -2 {
		i := 0
		for sqlExca[i] != -1 || sqlExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; sqlExca[i] >= 0; i += 2 {
			tok := sqlExca[i]
			if tok < TOKSTART || sqlExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if sqlExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += sqlTokname(tok)
	}
	return res
}

func sqllex1(lex sqlLexer, lval *sqlSymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = sqlTok1[0]
		goto out
	}
	if char < len(sqlTok1) {
		token = sqlTok1[char]
		goto out
	}
	if char >= sqlPrivate {
		if char < sqlPrivate+len(sqlTok2) {
			token = sqlTok2[char-sqlPrivate]
			goto out
		}
	}
	for i := 0; i < len(sqlTok3); i += 2 {
		token = sqlTok3[i+0]
		if token == char {
			token = sqlTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = sqlTok2[1] /* unknown char */
	}
	if sqlDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", sqlTokname(token), uint(char))
	}
	return char, token
}

func sqlParse(sqllex sqlLexer) int {
	return sqlNewParser().Parse(sqllex)
}

func (sqlrcvr *sqlParserImpl) Parse(sqllex sqlLexer) int {
	var sqln int
	var sqlVAL sqlSymType
	var sqlDollar []sqlSymType
	_ = sqlDollar // silence set and not used
	sqlS := sqlrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	sqlstate := 0
	sqlrcvr.char = -1
	sqltoken := -1 // sqlrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		sqlstate = -1
		sqlrcvr.char = -1
		sqltoken = -1
	}()
	sqlp := -1
	goto sqlstack

ret0:
	return 0

ret1:
	return 1

sqlstack:
	/* put a state and value onto the stack */
	if sqlDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", sqlTokname(sqltoken), sqlStatname(sqlstate))
	}

	sqlp++
	if sqlp >= len(sqlS) {
		nyys := make([]sqlSymType, len(sqlS)*2)
		copy(nyys, sqlS)
		sqlS = nyys
	}
	sqlS[sqlp] = sqlVAL
	sqlS[sqlp].yys = sqlstate

sqlnewstate:
	sqln = sqlPact[sqlstate]
	if sqln <= sqlFlag {
		goto sqldefault /* simple state */
	}
	if sqlrcvr.char < 0 {
		sqlrcvr.char, sqltoken = sqllex1(sqllex, &sqlrcvr.lval)
	}
	sqln += sqltoken
	if sqln < 0 || sqln >= sqlLast {
		goto sqldefault
	}
	sqln = sqlAct[sqln]
	if sqlChk[sqln] == sqltoken { /* valid shift */
		sqlrcvr.char = -1
		sqltoken = -1
		sqlVAL = sqlrcvr.lval
		sqlstate = sqln
		if Errflag > 0 {
			Errflag--
		}
		goto sqlstack
	}

sqldefault:
	/* default state action */
	sqln = sqlDef[sqlstate]
	if sqln == -2 {
		if sqlrcvr.char < 0 {
			sqlrcvr.char, sqltoken = sqllex1(sqllex, &sqlrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if sqlExca[xi+0] == -1 && sqlExca[xi+1] == sqlstate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			sqln = sqlExca[xi+0]
			if sqln < 0 || sqln == sqltoken {
				break
			}
		}
		sqln = sqlExca[xi+1]
		if sqln < 0 {
			goto ret0
		}
	}
	if sqln == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			sqllex.Error(sqlErrorMessage(sqlstate, sqltoken))
			Nerrs++
			if sqlDebug >= 1 {
				__yyfmt__.Printf("%s", sqlStatname(sqlstate))
				__yyfmt__.Printf(" saw %s\n", sqlTokname(sqltoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for sqlp >= 0 {
				sqln = sqlPact[sqlS[sqlp].yys] + sqlErrCode
				if sqln >= 0 && sqln < sqlLast {
					sqlstate = sqlAct[sqln] /* simulate a shift of "error" */
					if sqlChk[sqlstate] == sqlErrCode {
						goto sqlstack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if sqlDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", sqlS[sqlp].yys)
				}
				sqlp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if sqlDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", sqlTokname(sqltoken))
			}
			if sqltoken == sqlEofCode {
				goto ret1
			}
			sqlrcvr.char = -1
			sqltoken = -1
			goto sqlnewstate /* try again in the same state */
		}
	}

	/* reduction by production sqln */
	if sqlDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", sqln, sqlStatname(sqlstate))
	}

	sqlnt := sqln
	sqlpt := sqlp
	_ = sqlpt // guard against "declared and not used"

	sqlp -= sqlR2[sqln]
	// sqlp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if sqlp+1 >= len(sqlS) {
		nyys := make([]sqlSymType, len(sqlS)*2)
		copy(nyys, sqlS)
		sqlS = nyys
	}
	sqlVAL = sqlS[sqlp+1]

	/* consult goto table to find next state */
	sqln = sqlR1[sqln]
	sqlg := sqlPgo[sqln]
	sqlj := sqlg + sqlS[sqlp].yys + 1

	if sqlj >= sqlLast {
		sqlstate = sqlAct[sqlg]
	} else {
		sqlstate = sqlAct[sqlj]
		if sqlChk[sqlstate] != -sqln {
			sqlstate = sqlAct[sqlg]
		}
	}
	// dummy call; replaced with literal code
	switch sqlnt {

	case 1:
		sqlDollar = sqlS[sqlpt-2 : sqlpt+1]
//line sql.y:146
		{
			parseResult = &extendedSelect{
				extended:       false,
				standardSelect: sqlDollar[1].slct}
		}
	case 2:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:151
		{
			parseResult = &extendedSelect{
				extended:       true,
				train:          true,
				standardSelect: sqlDollar[1].slct,
				trainClause:    sqlDollar[2].tran}
		}
	case 3:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:158
		{
			parseResult = &extendedSelect{
				extended:       true,
				train:          false,
				standardSelect: sqlDollar[1].slct,
				predictClause:  sqlDollar[2].infr}
		}
	case 4:
		sqlDollar = sqlS[sqlpt-2 : sqlpt+1]
//line sql.y:168
		{
			sqlVAL.slct.fields = sqlDollar[2].flds
		}
	case 5:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:169
		{
			sqlVAL.slct.tables = sqlDollar[3].tbls
		}
	case 6:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:170
		{
			sqlVAL.slct.limit = sqlDollar[3].val
		}
	case 7:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:171
		{
			sqlVAL.slct.where = sqlDollar[3].expr
		}
	case 8:
		sqlDollar = sqlS[sqlpt-8 : sqlpt+1]
//line sql.y:175
		{
			sqlVAL.tran.estimator = sqlDollar[2].val
			sqlVAL.tran.attrs = sqlDollar[4].atrs
			sqlVAL.tran.columns = sqlDollar[5].colc
			sqlVAL.tran.label = sqlDollar[6].labc
			sqlVAL.tran.save = sqlDollar[8].val
		}
	case 9:
		sqlDollar = sqlS[sqlpt-4 : sqlpt+1]
//line sql.y:185
		{
			sqlVAL.infr.into = sqlDollar[2].val
			sqlVAL.infr.model = sqlDollar[4].val
		}
	case 10:
		sqlDollar = sqlS[sqlpt-2 : sqlpt+1]
//line sql.y:192
		{
			sqlVAL.colc = map[string]exprlist{"feature_columns": sqlDollar[2].expl}
		}
	case 11:
		sqlDollar = sqlS[sqlpt-4 : sqlpt+1]
//line sql.y:193
		{
			sqlVAL.colc = map[string]exprlist{sqlDollar[4].val: sqlDollar[2].expl}
		}
	case 12:
		sqlDollar = sqlS[sqlpt-5 : sqlpt+1]
//line sql.y:194
		{
			sqlVAL.colc[sqlDollar[5].val] = sqlDollar[3].expl
		}
	case 13:
		sqlDollar = sqlS[sqlpt-1 : sqlpt+1]
//line sql.y:198
		{
			sqlVAL.flds = append(sqlVAL.flds, sqlDollar[1].val)
		}
	case 14:
		sqlDollar = sqlS[sqlpt-1 : sqlpt+1]
//line sql.y:199
		{
			sqlVAL.flds = append(sqlVAL.flds, sqlDollar[1].val)
		}
	case 15:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:200
		{
			sqlVAL.flds = append(sqlVAL.flds, sqlDollar[3].val)
		}
	case 16:
		sqlDollar = sqlS[sqlpt-1 : sqlpt+1]
//line sql.y:204
		{
			sqlVAL.expr = atomic(IDENT, "*")
		}
	case 17:
		sqlDollar = sqlS[sqlpt-1 : sqlpt+1]
//line sql.y:205
		{
			sqlVAL.expr = atomic(IDENT, sqlDollar[1].val)
		}
	case 18:
		sqlDollar = sqlS[sqlpt-1 : sqlpt+1]
//line sql.y:206
		{
			sqlVAL.expr = sqlDollar[1].expr
		}
	case 19:
		sqlDollar = sqlS[sqlpt-1 : sqlpt+1]
//line sql.y:210
		{
			sqlVAL.expl = exprlist{sqlDollar[1].expr}
		}
	case 20:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:211
		{
			sqlVAL.expl = append(sqlDollar[1].expl, sqlDollar[3].expr)
		}
	case 21:
		sqlDollar = sqlS[sqlpt-2 : sqlpt+1]
//line sql.y:215
		{
			sqlVAL.labc = sqlDollar[2].val
		}
	case 22:
		sqlDollar = sqlS[sqlpt-2 : sqlpt+1]
//line sql.y:216
		{
			sqlVAL.labc = sqlDollar[2].val
		}
	case 23:
		sqlDollar = sqlS[sqlpt-1 : sqlpt+1]
//line sql.y:220
		{
			sqlVAL.tbls = []string{sqlDollar[1].val}
		}
	case 24:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:221
		{
			sqlVAL.tbls = append(sqlDollar[1].tbls, sqlDollar[3].val)
		}
	case 25:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:225
		{
			sqlVAL.atrs = attrs{sqlDollar[1].val: sqlDollar[3].expr}
		}
	case 26:
		sqlDollar = sqlS[sqlpt-1 : sqlpt+1]
//line sql.y:229
		{
			sqlVAL.atrs = sqlDollar[1].atrs
		}
	case 27:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:230
		{
			sqlVAL.atrs = attrsUnion(sqlDollar[1].atrs, sqlDollar[3].atrs)
		}
	case 28:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:234
		{
			sqlVAL.expr = funcall(sqlDollar[1].val, nil)
		}
	case 29:
		sqlDollar = sqlS[sqlpt-4 : sqlpt+1]
//line sql.y:235
		{
			sqlVAL.expr = funcall(sqlDollar[1].val, sqlDollar[3].expl)
		}
	case 30:
		sqlDollar = sqlS[sqlpt-1 : sqlpt+1]
//line sql.y:239
		{
			sqlVAL.expl = exprlist{sqlDollar[1].expr}
		}
	case 31:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:240
		{
			sqlVAL.expl = append(sqlDollar[1].expl, sqlDollar[3].expr)
		}
	case 32:
		sqlDollar = sqlS[sqlpt-2 : sqlpt+1]
//line sql.y:244
		{
			sqlVAL.expl = nil
		}
	case 33:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:245
		{
			sqlVAL.expl = sqlDollar[2].expl
		}
	case 34:
		sqlDollar = sqlS[sqlpt-1 : sqlpt+1]
//line sql.y:249
		{
			sqlVAL.expr = atomic(NUMBER, sqlDollar[1].val)
		}
	case 35:
		sqlDollar = sqlS[sqlpt-1 : sqlpt+1]
//line sql.y:250
		{
			sqlVAL.expr = atomic(IDENT, sqlDollar[1].val)
		}
	case 36:
		sqlDollar = sqlS[sqlpt-1 : sqlpt+1]
//line sql.y:251
		{
			sqlVAL.expr = atomic(STRING, sqlDollar[1].val)
		}
	case 37:
		sqlDollar = sqlS[sqlpt-1 : sqlpt+1]
//line sql.y:252
		{
			sqlVAL.expr = variadic('[', "square", sqlDollar[1].expl)
		}
	case 38:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:253
		{
			sqlVAL.expr = unary('(', "paren", sqlDollar[2].expr)
		}
	case 39:
		sqlDollar = sqlS[sqlpt-1 : sqlpt+1]
//line sql.y:254
		{
			sqlVAL.expr = sqlDollar[1].expr
		}
	case 40:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:255
		{
			sqlVAL.expr = binary('+', sqlDollar[1].expr, sqlDollar[2].val, sqlDollar[3].expr)
		}
	case 41:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:256
		{
			sqlVAL.expr = binary('-', sqlDollar[1].expr, sqlDollar[2].val, sqlDollar[3].expr)
		}
	case 42:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:257
		{
			sqlVAL.expr = binary('*', sqlDollar[1].expr, sqlDollar[2].val, sqlDollar[3].expr)
		}
	case 43:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:258
		{
			sqlVAL.expr = binary('/', sqlDollar[1].expr, sqlDollar[2].val, sqlDollar[3].expr)
		}
	case 44:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:259
		{
			sqlVAL.expr = binary('%', sqlDollar[1].expr, sqlDollar[2].val, sqlDollar[3].expr)
		}
	case 45:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:260
		{
			sqlVAL.expr = binary('=', sqlDollar[1].expr, sqlDollar[2].val, sqlDollar[3].expr)
		}
	case 46:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:261
		{
			sqlVAL.expr = binary('<', sqlDollar[1].expr, sqlDollar[2].val, sqlDollar[3].expr)
		}
	case 47:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:262
		{
			sqlVAL.expr = binary('>', sqlDollar[1].expr, sqlDollar[2].val, sqlDollar[3].expr)
		}
	case 48:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:263
		{
			sqlVAL.expr = binary(LE, sqlDollar[1].expr, sqlDollar[2].val, sqlDollar[3].expr)
		}
	case 49:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:264
		{
			sqlVAL.expr = binary(GE, sqlDollar[1].expr, sqlDollar[2].val, sqlDollar[3].expr)
		}
	case 50:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:265
		{
			sqlVAL.expr = binary(AND, sqlDollar[1].expr, sqlDollar[2].val, sqlDollar[3].expr)
		}
	case 51:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:266
		{
			sqlVAL.expr = binary(OR, sqlDollar[1].expr, sqlDollar[2].val, sqlDollar[3].expr)
		}
	case 52:
		sqlDollar = sqlS[sqlpt-2 : sqlpt+1]
//line sql.y:267
		{
			sqlVAL.expr = unary(NOT, sqlDollar[1].val, sqlDollar[2].expr)
		}
	case 53:
		sqlDollar = sqlS[sqlpt-2 : sqlpt+1]
//line sql.y:268
		{
			sqlVAL.expr = unary('-', sqlDollar[1].val, sqlDollar[2].expr)
		}
	}
	goto sqlstack /* stack new state and value */
}
