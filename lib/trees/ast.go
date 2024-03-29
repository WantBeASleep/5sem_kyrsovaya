package trees

import (
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"

	mt "lib/matrixes"
)

type ASTNode interface {
	GetMatrixSize() mt.MatrixSize
	GetCountOp() int
}

type MatrixLeaf struct {
	MatrixName string
	Size       mt.MatrixSize
}

func (m *MatrixLeaf) GetMatrixSize() mt.MatrixSize {
	return m.Size
}

func (m *MatrixLeaf) GetCountOp() int {
	return m.Size.Rows * m.Size.Cols
}

type BinaryOp struct {
	Op    token.Token
	Left  ASTNode
	Right ASTNode

	Size                   mt.MatrixSize
	SubTreeCountOperations int
}

func (b *BinaryOp) GetMatrixSize() mt.MatrixSize {
	return b.Size
}

func (b *BinaryOp) GetCountOp() int {
	return b.SubTreeCountOperations
}

func ParseExpr(expr string) ASTNode {
	fset := token.NewFileSet()
	defAst, _ := parser.ParseExprFrom(fset, "", expr, 0)
	tree := parseGoAstWithoutSize(defAst)
	return tree
}

func parseGoAstWithoutSize(n ast.Node) ASTNode {
	switch x := n.(type) {
	case *ast.ParenExpr:
		return parseGoAstWithoutSize(x.X)

	case *ast.BinaryExpr:
		left := parseGoAstWithoutSize(x.X)
		right := parseGoAstWithoutSize(x.Y)

		newNode := BinaryOp{
			Op:    x.Op,
			Left:  left,
			Right: right,
		}

		return &newNode

	case *ast.Ident:
		newNode := MatrixLeaf{
			MatrixName: x.Name,
		}
		return &newNode
	}
	return nil
}

func UpdateTreeStats(node ASTNode, data map[string]mt.Matrix) {
	var dfs func(nd ASTNode)
	dfs = func(nd ASTNode) {
		switch x := nd.(type) {
		case *BinaryOp:
			dfs(x.Left)
			dfs(x.Right)

			switch x.Op {
			case token.ADD, token.SUB:
				x.Size = x.Left.GetMatrixSize()

				opWeight := x.Size.Rows * x.Size.Cols
				x.SubTreeCountOperations = x.Left.GetCountOp() + x.Right.GetCountOp() + opWeight

			case token.MUL:
				x.Size = mt.MatrixSize{
					Rows: x.Left.GetMatrixSize().Rows,
					Cols: x.Right.GetMatrixSize().Cols,
				}

				opWeight := x.Size.Rows * x.Size.Cols * (x.Left.GetMatrixSize().Cols * x.Right.GetMatrixSize().Rows)
				x.SubTreeCountOperations = x.Left.GetCountOp() + x.Right.GetCountOp() + opWeight
			}

		case *MatrixLeaf:
			x.Size = data[x.MatrixName].Size
		}
	}
	dfs(node)
}

func GetLeafsNames(root ASTNode) map[string]bool {
	answr := map[string]bool{}

	var dfs func(nd ASTNode)
	dfs = func(nd ASTNode) {
		switch x := nd.(type) {
		case *BinaryOp:
			dfs(x.Left)
			dfs(x.Right)

		case *MatrixLeaf:
			answr[x.MatrixName] = true
		}
	}
	dfs(root)

	return answr
}


func UpparseJson(tree json.RawMessage) ASTNode {
	var dfs func(json.RawMessage) ASTNode
	dfs = func(node json.RawMessage) ASTNode {
		var x map[string]json.RawMessage
		err := json.Unmarshal(node, &x)
		if err != nil {
			panic("ошибка при распарсе на воркере(в dfs)")
		}

		if _, isLeaf := x["MatrixName"]; isLeaf {
			var newLeaf MatrixLeaf
			json.Unmarshal(x["MatrixName"], &newLeaf.MatrixName)
			json.Unmarshal(x["Size"], &newLeaf.Size)
			return &newLeaf
		} else {
			left := dfs(x["Left"])
			right := dfs(x["Left"])

			var newBin BinaryOp
			json.Unmarshal(x["Op"], &newBin.Op)
			newBin.Left = left
			newBin.Right = right
			json.Unmarshal(x["Size"], &newBin.Size)
			json.Unmarshal(x["SubTreeCountOperations"], &newBin.SubTreeCountOperations)
			return &newBin
		}
	}
	return dfs(tree)
}