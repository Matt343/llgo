// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ast

import "fmt"

// A Visitor's Visit method is invoked for each node encountered by Walk.
// If the result visitor w is not nil, Walk visits each of the children
// of node with the visitor w, followed by a call of w.Visit(nil).
type Visitor interface {
	Visit(node Node) (w Visitor)
}

// Helper functions for common node lists. They may be empty.

func walkIdentList(v Visitor, list []*Ident) {
	for _, x := range list {
		Walk(v, x)
	}
}

func walkExprList(v Visitor, list []Expr) {
	for _, x := range list {
		Walk(v, x)
	}
}

func walkStmtList(v Visitor, list []Stmt) {
	for _, x := range list {
		Walk(v, x)
	}
}

func walkDeclList(v Visitor, list []Decl) {
	for _, x := range list {
		Walk(v, x)
	}
}

// TODO(gri): Investigate if providing a closure to Walk leads to
//            simpler use (and may help eliminate Inspect in turn).

// Walk traverses an AST in depth-first order: It starts by calling
// v.Visit(node); node must not be nil. If the visitor w returned by
// v.Visit(node) is not nil, Walk is invoked recursively with visitor
// w for each of the non-nil children of node, followed by a call of
// w.Visit(nil).
//
func Walk(v Visitor, node Node) {
	if v = v.Visit(node); v == nil {
		return
	}

	// walk children
	// (the order of the cases matches the order
	// of the corresponding node types in ast.go)
	switch n := node.(type) {
	// Comments and fields
	case *Comment:
		// nothing to do

	case *CommentGroup:
		for _, c := range n.List {
			Walk(v, c)
		}

	case *Field:
		if n.Doc != nil {
			Walk(v, n.Doc)
		}
		walkIdentList(v, n.Names)
		Walk(v, n.Type)
		if n.Tag != nil {
			Walk(v, n.Tag)
		}
		if n.Comment != nil {
			Walk(v, n.Comment)
		}

	case *FieldList:
		for _, f := range n.List {
			Walk(v, f)
		}

	case *TypeParameter:
		if n.Doc != nil {
			Walk(v, n.Doc)
		}
		walkIdentList(v, n.Names)
		Walk(v, n.TypeBound)
		if n.Tag != nil {
			Walk(v, n.Tag)
		}
		if n.Comment != nil {
			Walk(v, n.Comment)
		}

	case *TypeParameterList:
		for _, t := range n.List {
			Walk(v, t)
		}

	// Expressions
	case *BadExpr, *Ident, *BasicLit:
		// nothing to do

	case *Ellipsis:
		if n.Elt != nil {
			Walk(v, n.Elt)
		}

	case *FuncLit:
		Walk(v, n.Type)
		Walk(v, n.Body)

	case *CompositeLit:
		if n.Type != nil {
			Walk(v, n.Type)
		}
		walkExprList(v, n.Elts)

	case *ParenExpr:
		Walk(v, n.X)

	case *SelectorExpr:
		Walk(v, n.X)
		Walk(v, n.Sel)

	case *IndexExpr:
		Walk(v, n.X)
		Walk(v, n.Index)

	case *SliceExpr:
		Walk(v, n.X)
		if n.Low != nil {
			Walk(v, n.Low)
		}
		if n.High != nil {
			Walk(v, n.High)
		}
		if n.Max != nil {
			Walk(v, n.Max)
		}

	case *TypeAssertExpr:
		Walk(v, n.X)
		if n.Type != nil {
			Walk(v, n.Type)
		}

	case *CallExpr:
		Walk(v, n.Fun)
		walkExprList(v, n.Args)

	case *StarExpr:
		Walk(v, n.X)

	case *UnaryExpr:
		Walk(v, n.X)

	case *BinaryExpr:
		Walk(v, n.X)
		Walk(v, n.Y)

	case *KeyValueExpr:
		Walk(v, n.Key)
		Walk(v, n.Value)

	// Types
	case *ArrayType:
		if n.Len != nil {
			Walk(v, n.Len)
		}
		Walk(v, n.Elt)

	case *StructType:
		if n.TypeParams != nil {
			Walk(v, n.TypeParams)
		}
		Walk(v, n.Fields)

	case *FuncType:
		if n.TypeParams != nil {
			Walk(v, n.TypeParams)
		}
		if n.Params != nil {
			Walk(v, n.Params)
		}
		if n.Results != nil {
			Walk(v, n.Results)
		}

	case *InterfaceType:
		if n.TypeParams != nil {
			Walk(v, n.TypeParams)
		}
		Walk(v, n.Methods)

	case *MapType:
		Walk(v, n.Key)
		Walk(v, n.Value)

	case *ChanType:
		Walk(v, n.Value)

	case *GenericType:
		Walk(v, n.Type)
		walkExprList(v, n.TypeParameters)

	// Statements
	case *BadStmt:
		// nothing to do

	case *DeclStmt:
		Walk(v, n.Decl)

	case *EmptyStmt:
		// nothing to do

	case *LabeledStmt:
		Walk(v, n.Label)
		Walk(v, n.Stmt)

	case *ExprStmt:
		Walk(v, n.X)

	case *SendStmt:
		Walk(v, n.Chan)
		Walk(v, n.Value)

	case *IncDecStmt:
		Walk(v, n.X)

	case *AssignStmt:
		walkExprList(v, n.Lhs)
		walkExprList(v, n.Rhs)

	case *GoStmt:
		Walk(v, n.Call)

	case *DeferStmt:
		Walk(v, n.Call)

	case *ReturnStmt:
		walkExprList(v, n.Results)

	case *BranchStmt:
		if n.Label != nil {
			Walk(v, n.Label)
		}

	case *BlockStmt:
		walkStmtList(v, n.List)

	case *IfStmt:
		if n.Init != nil {
			Walk(v, n.Init)
		}
		Walk(v, n.Cond)
		Walk(v, n.Body)
		if n.Else != nil {
			Walk(v, n.Else)
		}

	case *CaseClause:
		walkExprList(v, n.List)
		walkStmtList(v, n.Body)

	case *SwitchStmt:
		if n.Init != nil {
			Walk(v, n.Init)
		}
		if n.Tag != nil {
			Walk(v, n.Tag)
		}
		Walk(v, n.Body)

	case *TypeSwitchStmt:
		if n.Init != nil {
			Walk(v, n.Init)
		}
		Walk(v, n.Assign)
		Walk(v, n.Body)

	case *CommClause:
		if n.Comm != nil {
			Walk(v, n.Comm)
		}
		walkStmtList(v, n.Body)

	case *SelectStmt:
		Walk(v, n.Body)

	case *ForStmt:
		if n.Init != nil {
			Walk(v, n.Init)
		}
		if n.Cond != nil {
			Walk(v, n.Cond)
		}
		if n.Post != nil {
			Walk(v, n.Post)
		}
		Walk(v, n.Body)

	case *RangeStmt:
		if n.Key != nil {
			Walk(v, n.Key)
		}
		if n.Value != nil {
			Walk(v, n.Value)
		}
		Walk(v, n.X)
		Walk(v, n.Body)

	// Declarations
	case *ImportSpec:
		if n.Doc != nil {
			Walk(v, n.Doc)
		}
		if n.Name != nil {
			Walk(v, n.Name)
		}
		Walk(v, n.Path)
		if n.Comment != nil {
			Walk(v, n.Comment)
		}

	case *ValueSpec:
		if n.Doc != nil {
			Walk(v, n.Doc)
		}
		walkIdentList(v, n.Names)
		if n.Type != nil {
			Walk(v, n.Type)
		}
		walkExprList(v, n.Values)
		if n.Comment != nil {
			Walk(v, n.Comment)
		}

	case *TypeSpec:
		if n.Doc != nil {
			Walk(v, n.Doc)
		}
		Walk(v, n.Name)
		Walk(v, n.Type)
		if n.Comment != nil {
			Walk(v, n.Comment)
		}

	case *BadDecl:
		// nothing to do

	case *GenDecl:
		if n.Doc != nil {
			Walk(v, n.Doc)
		}
		for _, s := range n.Specs {
			Walk(v, s)
		}

	case *FuncDecl:
		if n.Doc != nil {
			Walk(v, n.Doc)
		}
		if n.Recv != nil {
			Walk(v, n.Recv)
		}
		Walk(v, n.Name)
		Walk(v, n.Type)
		if n.Body != nil {
			Walk(v, n.Body)
		}

	// Files and packages
	case *File:
		if n.Doc != nil {
			Walk(v, n.Doc)
		}
		Walk(v, n.Name)
		walkDeclList(v, n.Decls)
		// don't walk n.Comments - they have been
		// visited already through the individual
		// nodes

	case *Package:
		for _, f := range n.Files {
			Walk(v, f)
		}

	default:
		panic(fmt.Sprintf("ast.Walk: unexpected node type %T", n))
	}

	v.Visit(nil)
}

type inspector func(Node) bool

func (f inspector) Visit(node Node) Visitor {
	if f(node) {
		return f
	}
	return nil
}

// Inspect traverses an AST in depth-first order: It starts by calling
// f(node); node must not be nil. If f returns true, Inspect invokes f
// recursively for each of the non-nil children of node, followed by a
// call of f(nil).
//
func Inspect(node Node, f func(Node) bool) {
	Walk(inspector(f), node)
}



type Transformer interface {
	TransformNode(old, input Node) (output Node)
	TransformExpr(old, input Expr) (output Expr)
	TransformDecl(old, input Decl) (output Decl)
	TransformStmt(old, input Stmt) (output []Stmt)
	TransformSpec(old, input Spec) (output Spec)
}

func walkTransformIdentList(v Transformer, list []*Ident) (output []*Ident){
	for _, x := range list {
		output = append(output, WalkTransformExpr(v, x).(*Ident))
	}
	return
}

func walkTransformExprList(v Transformer, list []Expr) (output []Expr) {
	for _, x := range list {
		output = append(output, WalkTransformExpr(v, x))
	}
	return
}

func walkTransformStmtList(v Transformer, list []Stmt) (output []Stmt) {
	for _, x := range list {
		output = append(output, WalkTransformStmt(v, x)...)
	}
	return
}

func walkTransformDeclList(v Transformer, list []Decl) (output []Decl) {
	for _, x := range list {
		output = append(output, WalkTransformDecl(v, x))
	}
	return
}

func walkTransformSpecList(v Transformer, list []Spec) (output []Spec) {
	for _, x := range list {
		output = append(output, WalkTransformSpec(v, x).(Spec))
	}
	return
}

func walkTransformInitStmt(v Transformer, init Stmt) (output []Stmt, newInit Stmt) {
	newInitList := WalkTransformStmt(v, init)
	length := len(newInitList)
	if length == 1 {
		newInit = newInitList[0]
	} else if length > 1 {
		output = newInitList[0:length - 1]
		newInit = newInitList[length - 1]
	}
	return
}

func WalkTransformExpr(v Transformer, node Expr) Expr {
	// walk children
	// (the order of the cases matches the order
	// of the corresponding node types in ast.go)
	switch n := node.(type) {

	case *FuncLit:
		return v.TransformExpr(n, &FuncLit{n.Type, WalkTransformStmt(v, n.Body)[0].(*BlockStmt)})

	case *CompositeLit:
		newElts := walkTransformExprList(v, n.Elts)
		return v.TransformExpr(n, &CompositeLit{n.Type, n.Lbrace, newElts, n.Rbrace})

	case *ParenExpr:
		return v.TransformExpr(n, &ParenExpr{n.Lparen, WalkTransformExpr(v, n.X), n.Rparen})

	case *SelectorExpr:
		newX := WalkTransformExpr(v, n.X)
		newSel := WalkTransformExpr(v, n.Sel).(*Ident)
		return v.TransformExpr(n, &SelectorExpr{newX, newSel})

	case *IndexExpr:
		newX := WalkTransformExpr(v, n.X)
		newIndex := WalkTransformExpr(v, n.Index)
		return v.TransformExpr(n, &IndexExpr{newX, n.Lbrack, newIndex, n.Rbrack})

	// case *SliceExpr:
	// 	Walk(v, n.X)
	// 	if n.Low != nil {
	// 		Walk(v, n.Low)
	// 	}
	// 	if n.High != nil {
	// 		Walk(v, n.High)
	// 	}
	// 	if n.Max != nil {
	// 		Walk(v, n.Max)
	// 	}

	// case *TypeAssertExpr:
	// 	Walk(v, n.X)
	// 	if n.Type != nil {
	// 		Walk(v, n.Type)
	// 	}

	case *CallExpr:
		newFun := WalkTransformExpr(v, n.Fun)
		newArgs := walkTransformExprList(v, n.Args)
		return v.TransformExpr(n, &CallExpr{newFun, n.Lbrack, n.TypeArgs, n.Rbrack, n.Lparen, newArgs, n.Ellipsis, n.Rparen})

	case *StarExpr:
		return v.TransformExpr(n, &StarExpr{n.Star, WalkTransformExpr(v, n.X)})

	case *UnaryExpr:
		return v.TransformExpr(n, &UnaryExpr{n.OpPos, n.Op, WalkTransformExpr(v, n.X)})

	case *BinaryExpr:
		newX := WalkTransformExpr(v, n.X)
		newY := WalkTransformExpr(v, n.Y)
		return v.TransformExpr(n, &BinaryExpr{newX, n.OpPos, n.Op, newY})

	case *KeyValueExpr:
		newKey := WalkTransformExpr(v, n.Key)
		newValue := WalkTransformExpr(v, n.Value)
		return v.TransformExpr(n, &KeyValueExpr{newKey, n.Colon, newValue})

	// Types
	case *ArrayType, *StructType, *FuncType, *InterfaceType, *MapType, *ChanType:
		return v.TransformExpr(n, n)

	default:
		return v.TransformExpr(n, n)
	}
}

func WalkTransformStmt(v Transformer, node Stmt) (output []Stmt) {
	switch n := node.(type) {
	case *DeclStmt:
		return v.TransformStmt(n, &DeclStmt{WalkTransformDecl(v, n.Decl)})

	case *LabeledStmt:
		newLabel := WalkTransformExpr(v, n.Label).(*Ident)
		var newStmt Stmt
		newStmts := WalkTransformStmt(v, n.Stmt)
		lenNewStmts := len(newStmts)
		if lenNewStmts == 1 {
			newStmt = newStmts[0]
		} else if lenNewStmts > 1 {
			output = newStmts[0:lenNewStmts-1]
			newStmt = newStmts[lenNewStmts-1]
		}

		output = append(output, v.TransformStmt(n, &LabeledStmt{newLabel, n.Colon, newStmt})...)
		return

	case *ExprStmt:
		return v.TransformStmt(n, &ExprStmt{WalkTransformExpr(v, n.X)})

	case *SendStmt:
		newChan := WalkTransformExpr(v, n.Chan)
		newValue := WalkTransformExpr(v, n.Value)
		return v.TransformStmt(n, &SendStmt{newChan, n.Arrow, newValue})

	case *IncDecStmt:
		return v.TransformStmt(n, &IncDecStmt{WalkTransformExpr(v, n.X), n.TokPos, n.Tok})

	case *AssignStmt:
		newLhs := walkTransformExprList(v, n.Lhs)
		newRhs := walkTransformExprList(v, n.Rhs)
		return v.TransformStmt(n, &AssignStmt{newLhs, n.TokPos, n.Tok, newRhs})

	case *GoStmt:
		return v.TransformStmt(n, &GoStmt{n.Go, WalkTransformExpr(v, n.Call).(*CallExpr)})

	case *DeferStmt:
		return v.TransformStmt(n, &DeferStmt{n.Defer, WalkTransformExpr(v, n.Call).(*CallExpr)})

	case *ReturnStmt:
		return v.TransformStmt(n, &ReturnStmt{n.Return, walkTransformExprList(v, n.Results)})

	case *BlockStmt:
		return v.TransformStmt(n, &BlockStmt{n.Lbrace, walkTransformStmtList(v, n.List), n.Rbrace})

	case *IfStmt:
		output, newInit := walkTransformInitStmt(v, n.Init)

		newCond := WalkTransformExpr(v, n.Cond)

		newBody := WalkTransformStmt(v, n.Body)[0].(*BlockStmt)

		newElseList := WalkTransformStmt(v, n.Else)
		var newElse Stmt
		if len(newElseList) == 1 {
			newElse = newElseList[0]
		} else if len(newElseList) > 1 {
			newElse = &BlockStmt{List: newElseList}
		}

		output = append(output, v.TransformStmt(n, &IfStmt{n.If, newInit, newCond, newBody, newElse})...)
		return output

	case *CaseClause:
		newList := walkTransformExprList(v, n.List)
		newBody := walkTransformStmtList(v, n.Body)
		return v.TransformStmt(n, &CaseClause{n.Case, newList, n.Colon, newBody})

	case *SwitchStmt:
		output, newInit := walkTransformInitStmt(v, n.Init)
		newBody := WalkTransformStmt(v, n.Body)[0].(*BlockStmt)
		output = append(output, v.TransformStmt(n, &SwitchStmt{n.Switch, newInit, n.Tag, newBody})...)
		return output

	case *TypeSwitchStmt:
		output, newInit := walkTransformInitStmt(v, n.Init)
		newAssign := WalkTransformStmt(v, n.Assign)[0]
		newBody := WalkTransformStmt(v, n.Body)[0].(*BlockStmt)
		output = append(output, v.TransformStmt(n, &TypeSwitchStmt{n.Switch, newInit, newAssign, newBody})...)
		return output

	case *CommClause:
		var newComm Stmt
		if n.Comm != nil {
			newComm = WalkTransformStmt(v, n.Comm)[0]
		}
		newBody := walkTransformStmtList(v, n.Body)
		return v.TransformStmt(n, &CommClause{n.Case, newComm, n.Colon, newBody})

	case *SelectStmt:
		return v.TransformStmt(n, &SelectStmt{n.Select, WalkTransformStmt(v, n.Body)[0].(*BlockStmt)})

	case *ForStmt:
		output, newInit := walkTransformInitStmt(v, n.Init)
		var newCond Expr
		if n.Cond != nil {
			newCond = WalkTransformExpr(v, n.Cond)
		}
		newBody := WalkTransformStmt(v, n.Body)[0].(*BlockStmt)

		var newPost Stmt
		newPostList := WalkTransformStmt(v, n.Post)
		lenPostList := len(newPostList)
		if lenPostList == 1 {
			newPost = newPostList[0]
		} else if lenPostList > 1 {
			newBody.List = append(newBody.List, newPostList[0:lenPostList - 1]...)
			newPost = newPostList[lenPostList - 1]
		}

		output = append(output, v.TransformStmt(n, &ForStmt{n.For, newInit, newCond, newPost, newBody})...)
		return output

	case *RangeStmt:
		var newKey, newValue Expr
		if n.Key != nil {
			newKey = WalkTransformExpr(v, n.Key)
		}
		if n.Value != nil {
			newValue = WalkTransformExpr(v, n.Value)
		}
		newX := WalkTransformExpr(v, n.X)
		newBody := WalkTransformStmt(v, n.Body)[0].(*BlockStmt)
		return v.TransformStmt(n, &RangeStmt{n.For, newKey, newValue, n.TokPos, n.Tok, newX, newBody})
	
	default:
		return v.TransformStmt(n, n)
	}
}

func WalkTransformSpec(v Transformer, node Spec) Spec {
	switch n := node.(type) {
	case *ValueSpec:
		newNames := walkTransformIdentList(v, n.Names)
		var newType Expr
		if n.Type != nil {
			newType = WalkTransformExpr(v, n.Type)
		}
		newValues := walkTransformExprList(v, n.Values)
		return v.TransformSpec(n, &ValueSpec{n.Doc, newNames, newType, newValues, n.Comment})

	case *TypeSpec:
		newType := WalkTransformExpr(v, n.Type)
		return v.TransformSpec(n, &TypeSpec{n.Doc, n.Name, newType, n.Comment})
		
	default:
		return v.TransformSpec(n, n)
	}
}

func WalkTransformDecl(v Transformer, node Decl) Decl {
	switch n := node.(type) {
	case *GenDecl:
		newSpecs := walkTransformSpecList(v, n.Specs)
		return v.TransformDecl(n, &GenDecl{n.Doc, n.TokPos, n.Tok, n.Lparen, newSpecs, n.Rparen})

	case *FuncDecl:
		var newRecv *FieldList
		if n.Recv != nil {
			newRecv = WalkTransform(v, n.Recv).(*FieldList)
		}
		newType := WalkTransformExpr(v, n.Type).(*FuncType)
		var newBody *BlockStmt
		if n.Body != nil {
			newBody = WalkTransformStmt(v, n.Body)[0].(*BlockStmt)
		}
		return v.TransformDecl(n, &FuncDecl{n.Doc, newRecv, n.Name, newType, newBody})

	default:
		return v.TransformDecl(n, n)
	}
}

func WalkTransform(v Transformer, node Node) Node {
	switch n := node.(type) {
	// Files and packages
	case *File:
		newDecls := walkTransformDeclList(v, n.Decls)
		return v.TransformNode(n, &File{
			n.Doc,
			n.Package,
			n.Name,
			newDecls,
			n.Scope,
			n.Imports,
			n.Unresolved,
			n.Comments,
		})

	case *Package:
		newFiles := make(map[string]*File, 0)
		for name, f := range n.Files {
			newFiles[name] = WalkTransform(v, f).(*File)
		}
		return v.TransformNode(n, &Package{n.Name, n.Scope, n.Imports, newFiles})

	case Expr:
		return WalkTransformExpr(v, n)

	case Stmt:
		return WalkTransformStmt(v, n)[0]

	case Decl:
		return WalkTransformDecl(v, n)

	case Spec:
		return WalkTransformSpec(v, n)

	case *Field:
		newNames := walkTransformIdentList(v, n.Names)
		newType := WalkTransformExpr(v, n.Type)
		return v.TransformNode(n, &Field{n.Doc, newNames, newType, n.Tag, n.Comment})

	case *FieldList:
		newList := make([]*Field, 0)
		for _, f := range n.List {
			newList = append(newList, WalkTransform(v, f).(*Field))
		}
		return v.TransformNode(n, &FieldList{n.Opening, newList, n.Closing})

	default:
		return n
	}
}
