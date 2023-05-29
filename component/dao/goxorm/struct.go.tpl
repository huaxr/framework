package models

{{$ilen := len .Imports}}
import (
        {{range .Imports}}"{{.}}"{{end}}
       )

{{range .Tables}}
{{$tb := Mapper .Name}}
{{$table := .}}
{{$dao := printf "%sDao" $tb}}

type {{$tb}} struct {
    {{range .ColumnsSeq}}{{$col := $table.GetColumn .}} {{Mapper $col.Name}}    {{Type $col}} {{Tag $table $col}}
    {{end}}
}

{{end}}
