package models

import "sort"

type DataFrameProgram struct {
	Programs []*Program `json:"programs,omitempty"`
}

func (df *DataFrameProgram) Performers() []string {
	// パファーマー名のスライス生成
	s := make([]string, len(df.Programs))
	for i, program := range df.Programs {
		s[i] = program.Performer
	}

	// パフォーマーの重複を削除
	m := make(map[string]bool)
	performers := make([]string, 0)
	for _, performer := range s {
		_, is := m[performer]
		if !is {
			m[performer] = true
			performers = append(performers, performer)
		}
	}

	// 名前順にソート
	sort.Strings(performers)
	return performers
}

func (df *DataFrameProgram) Vol() []string {
	s := make([]string, len(df.Programs))
	for i, program := range df.Programs {
		s[i] = program.Vol
	}

	m := make(map[string]bool)
	vols := make([]string, 0)
	for _, vol := range s {
		_, is := m[vol]
		if !is {
			m[vol] = true
			vols = append(vols, vol)
		}
	}
	sort.Strings(vols)
	return vols
}
