package types

import "fmt"

func (check *Checker) substituteTypes(context, typ Type, argTyp Type, aliases *TypeAliases, seen map[Type]Type) Type {
	if typ == nil {
		return nil
	}
	if seen == nil {
		seen = make(map[Type]Type)
	}
	if seen[typ] != nil {
		return seen[typ]
	}
	seen[typ] = typ

	var sub Type

	switch t := typ.(type) {
	case *Array:
		var argElem Type
		if argArray, ok := argTyp.(*Array); ok {
			argElem = argArray.elem
		}
		sub = &Array{t.len, check.substituteTypes(context, t.elem, argElem, aliases, seen)}

	case *Slice:
		var argElem Type
		if argSlice, ok := argTyp.(*Slice); ok {
			argElem = argSlice.elem
		}
		sub = &Slice{check.substituteTypes(context, t.elem, argElem, aliases, seen)}

	case *Struct:
		var argFields []*Var
		var argTypeParams []*TypeName
		if argStruct, ok := argTyp.(*Struct); ok {
			argFields = argStruct.fields
			argTypeParams = argStruct.typeParams
		}
		sub = &Struct{
			check.substituteTypesVars(context, t.fields, argFields, aliases, seen),
			t.tags,
			t.offsets,
			check.substituteTypesTypeNames(context, t.typeParams, argTypeParams, aliases, seen),
		}

	case *Pointer:
		var argBase Type
		if argPointer, ok := argTyp.(*Pointer); ok {
			argBase = argPointer.base
		}
		sub = &Pointer{check.substituteTypes(context, t.base, argBase, aliases, seen)}

	case *Tuple:
		var argVars []*Var
		if argTuple, ok := argTyp.(*Tuple); ok {
			argVars = argTuple.vars
		}
		sub = &Tuple{check.substituteTypesVars(context, t.vars, argVars, aliases, seen)}

	case *Signature:
		var argParams *Tuple
		var argResults *Tuple
		var argTypeParams []*TypeName
		if argSignature, ok := argTyp.(*Signature); ok {
			argParams = argSignature.params
			argResults = argSignature.results
			argTypeParams = argSignature.typeParams
		}
		sub = &Signature{
			t.scope,
			t.recv,
			check.substituteTypesTuple(context, t.params, argParams, aliases, seen),
			check.substituteTypesTuple(context, t.results, argResults, aliases, seen),
			t.variadic,
			check.substituteTypesTypeNames(context, t.typeParams, argTypeParams, aliases, seen),
		}

	case *Interface:
		var argMethods []*Func
		var argEmbeds []*Named
		var argAllMethods []*Func
		if argInterface, ok := argTyp.(*Interface); ok {
			argMethods = argInterface.methods
			argEmbeds = argInterface.embeddeds
			argAllMethods = argInterface.allMethods
		}
		sub = &Interface{
			check.substituteTypesFuncs(context, t.methods, argMethods, aliases, seen),
			check.substituteTypesNameds(context, t.embeddeds, argEmbeds, aliases, seen),
			check.substituteTypesFuncs(context, t.allMethods, argAllMethods, aliases, seen),
			t.variance,
		}

	case *Map:
		var argKey Type
		var argElem Type
		if argMap, ok := argTyp.(*Map); ok {
			argKey = argMap.key
			argElem = argMap.elem
		}
		sub = &Map{
			check.substituteTypes(context, t.key, argKey, aliases, seen),
			check.substituteTypes(context, t.elem, argElem, aliases, seen),
		}

	case *Chan:
		var argElem Type
		if argChan, ok := argTyp.(*Chan); ok {
			argElem = argChan.elem
		}
		sub = &Chan{t.dir, check.substituteTypes(context, t.elem, argElem, aliases, seen)}

	case *Named:
		sub = check.substituteTypesNamed(context, t, argTyp, aliases, seen)

	default:
		sub = t
	}

	seen[typ] = sub
	return sub
}

func (check *Checker) substituteTypesNamed(context Type, old *Named, argTyp Type, aliases *TypeAliases, seen map[Type]Type) Type {
	if old == nil {
		return nil
	}
	if aliases != nil && old.obj != nil && old.context == context {
		if (*aliases)[old.obj] != nil {
			return (*aliases)[old.obj]
		} else if AssignableTo(argTyp, old) {
			(*aliases)[old.obj] = argTyp
			fmt.Printf("Infered %s -> %s\n", old.obj, argTyp)
			return argTyp
		} else {
			return old
		}
	} else {
		return old
	}
}

func (check *Checker) substituteTypesObject(context Type, old object, argObject object, aliases *TypeAliases, seen map[Type]Type) object {
	return object{old.parent, old.pos, old.pkg, old.name, check.substituteTypes(context, old.typ, argObject.typ, aliases, seen), old.order_}
}

func (check *Checker) substituteTypesVar(context Type, old *Var, argVar *Var, aliases *TypeAliases, seen map[Type]Type) *Var {
	if old == nil {
		return nil
	}
	var argObject object
	if argVar != nil {
		argObject = argVar.object
	}
	return &Var{check.substituteTypesObject(context, old.object, argObject, aliases, seen), old.anonymous, old.visited, old.isField, old.used}
}

func (check *Checker) substituteTypesFunc(context Type, old *Func, argFunc *Func, aliases *TypeAliases, seen map[Type]Type) *Func {
	if old == nil {
		return nil
	}
	var argObject object
	if argFunc != nil {
		argObject = argFunc.object
	}
	return &Func{check.substituteTypesObject(context, old.object, argObject, aliases, seen)}
}

func (check *Checker) substituteTypesTypeName(context Type, old *TypeName, argTypeName *TypeName, aliases *TypeAliases, seen map[Type]Type) *TypeName {
	if old == nil {
		return nil
	}
	var argObject object
	if argTypeName != nil {
		argObject = argTypeName.object
	}
	return &TypeName{check.substituteTypesObject(context, old.object, argObject, aliases, seen)}
}

func (check *Checker) substituteTypesNameds(context Type, old []*Named, argNameds []*Named, aliases *TypeAliases, seen map[Type]Type) []*Named {
	if old == nil {
		return nil
	}
	nameds := make([]*Named, len(old))
	for i, v := range old {
		var argNamed *Named
		if argNameds != nil && i < len(argNameds) {
			argNamed = argNameds[i]
		}
		nameds[i] = check.substituteTypesNamed(context, v, argNamed, aliases, seen).(*Named)
	}
	return nameds
}

func (check *Checker) substituteTypesVars(context Type, old []*Var, argVars []*Var, aliases *TypeAliases, seen map[Type]Type) []*Var {
	if old == nil {
		return nil
	}
	vars := make([]*Var, len(old))
	for i, v := range old {
		var argVar *Var
		if argVars != nil && i < len(argVars) {
			argVar = argVars[i]
		}
		vars[i] = check.substituteTypesVar(context, v, argVar, aliases, seen)
	}
	return vars
}

func (check *Checker) substituteTypesFuncs(context Type, old []*Func, argFuncs []*Func, aliases *TypeAliases, seen map[Type]Type) []*Func {
	if old == nil {
		return nil
	}
	funcs := make([]*Func, len(old))
	for i, f := range old {
		var argFunc *Func
		if argFuncs != nil && i < len(argFuncs) {
			argFunc = argFuncs[i]
		}
		funcs[i] = check.substituteTypesFunc(context, f, argFunc, aliases, seen)
	}
	return funcs
}

func (check *Checker) substituteTypesTypeNames(context Type, old []*TypeName, argTypeNames []*TypeName, aliases *TypeAliases, seen map[Type]Type) []*TypeName {
	if old == nil {
		return nil
	}
	names := make([]*TypeName, len(old))
	for i, t := range old {
		var argTypeName *TypeName
		if argTypeNames != nil && i < len(argTypeNames) {
			argTypeName = argTypeNames[i]
		}
		names[i] = check.substituteTypesTypeName(context, t, argTypeName, aliases, seen)
	}
	return names
}

func (check *Checker) substituteTypesTuple(context Type, old *Tuple, argTuple *Tuple, aliases *TypeAliases, seen map[Type]Type) *Tuple {
	if old == nil {
		return nil
	}
	var argVars []*Var
	if argTuple != nil {
		argVars = argTuple.vars
	}
	return &Tuple{check.substituteTypesVars(context, old.vars, argVars, aliases, seen)}
}
