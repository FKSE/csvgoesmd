package csvgoesmd

import (
    "testing"
    "io/ioutil"
)

func BenchmarkParseCsv(b *testing.B) {
    //csv content
    bytes, err := ioutil.ReadFile("./test.csv")
    //check for error
    if err != nil {
        panic(err)
    }
    content := string(bytes)
    b.ResetTimer()
    // run the Fib function b.N times
    for n := 0; n < b.N; n++ {
        ParseCsv(content, ";")
    }
}

func BenchmarkMaxLength(b *testing.B) {
    a := []string{"123", "sdfgerwgeg", "sdf43terge", "sad", "dsfsdfs"}
    for n := 0; n < b.N; n++ {
        MaxLength(a)
    }
}

func BenchmarkBuildMarkdown(b *testing.B) {
    table := [][]string{{"colA1", "colA2", "colA3", "colA4", "colA5"}, {"colB1", "colB2", "colB3", "colB4", "colB5"}, {"colC1", "colC2", "colC3", "colC4", "colC5"}}
    for n := 0; n < b.N; n++ {
        BuildMarkdown(table, false)
    }
}