package squ

import (
	"bytes"

	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	openapi "github.com/opensibyl/sibyl-go-client"
	"github.com/opensibyl/squ/indexer"
	"github.com/opensibyl/squ/object"
)

func renderGraph(outputPath string, curIndexer indexer.Indexer, diffMap map[string][]*object.FunctionWithState) error {
	buf := bytes.NewBuffer([]byte{})
	baseGraph := curIndexer.GetSibylCache().CallGraph.Graph
	err := draw.DOT(baseGraph, buf)
	PanicIfErr(err)

	g, err := graphviz.ParseBytes(buf.Bytes())
	g.SetRankDir(cgraph.LRRank)
	marker := &GraphMarker{}
	for _, funcList := range diffMap {
		for _, eachFunc := range funcList {
			if eachFunc.Reachable == false {
				// not covered
				vertexes := curIndexer.GetVertexesWithSignature(eachFunc.GetSignature())
				for _, eachVertex := range vertexes {
					node, err := g.Node(eachVertex)
					PanicIfErr(err)
					marker.MarkRed(node)
				}
				continue
			}
			// covered, good
			vertexes := curIndexer.GetVertexesWithSignature(eachFunc.GetSignature())
			for _, eachModifiedVertex := range vertexes {
				coveredModifiedNode, err := g.Node(eachModifiedVertex)
				PanicIfErr(err)
				marker.MarkGreen(coveredModifiedNode)

				// highlight the paths
				for _, eachCaseSignature := range eachFunc.ReachBy {
					caseVertexes := curIndexer.GetVertexesWithSignature(eachCaseSignature)
					for _, eachCaseVertex := range caseVertexes {
						paths, err := graph.ShortestPath(baseGraph, eachCaseVertex, eachModifiedVertex)
						PanicIfErr(err)
						// draw
						for _, eachPath := range paths {
							if eachPath == eachModifiedVertex {
								continue
							}
							node, err := g.Node(eachPath)
							PanicIfErr(err)
							marker.MarkYellow(node)
						}
					}
				}
			}
		}
	}
	if err != nil {
		return err
	}
	gviz := graphviz.New()
	err = gviz.RenderFilename(g, graphviz.PNG, outputPath)
	if err != nil {
		return err
	}
	return nil
}

type Output struct {
	DiffMap object.DiffFuncMap                     `json:"diff"`
	Cases   []*openapi.ObjectFunctionWithSignature `json:"cases"`
}

type GraphMarker struct {
}

func (g *GraphMarker) MarkRed(n *cgraph.Node) {
	n.SetStyle("filled")
	n.SetFillColor("red")
}

func (g *GraphMarker) MarkGreen(n *cgraph.Node) {
	n.SetStyle("filled")
	n.SetFillColor("green")
}

func (g *GraphMarker) MarkYellow(n *cgraph.Node) {
	n.SetStyle("filled")
	n.SetFillColor("greenyellow")
}
