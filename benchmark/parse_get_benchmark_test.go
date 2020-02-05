package jin

import (
	"github.com/json-iterator/go"
	"github.com/valyala/fastjson"
	"jin"
	"testing"
)

func BenchmarkJsoniteratorGetSmall(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		prs := jsoniter.Get(SmallFixture)
		prs.Get("uuid")
		prs.Get("tz")
		prs.Get("ua")
		prs.Get("st")
	}
}

func BenchmarkFastjsonGetSmall(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var p fastjson.Parser
		prs, _ := p.ParseBytes(SmallFixture)
		prs.Get("uuid")
		prs.Get("tz")
		prs.Get("ua")
		prs.Get("st")
	}
}

func BenchmarkJintParseGetSmall(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		prs, _ := jin.Parse(SmallFixture)
		prs.Get("uuid")
		prs.Get("tz")
		prs.Get("ua")
		prs.Get("st")
	}
}

func BenchmarkJsoniteratorGetMedium(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		prs := jsoniter.Get(MediumFixture)
		prs.Get("person", "name", "fullName")
		prs.Get("person", "github", "followers")
		prs.Get("company")
	}
}

func BenchmarkFastjsonGetMedium(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var p fastjson.Parser
		prs, _ := p.ParseBytes(MediumFixture)
		prs.Get("person", "name", "fullName")
		prs.Get("person", "github", "followers")
		prs.Get("company")
	}
}

func BenchmarkJintParseGetMedium(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		prs, _ := jin.Parse(MediumFixture)
		prs.Get("person", "name", "fullName")
		prs.Get("person", "github", "followers")
		prs.Get("company")
	}
}

func BenchmarkJsoniteratorGetLarge(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		prs := jsoniter.Get(LargeFixture)
		prs.Get("users", 0, "id")
		prs.Get("users", 31, "id")
		prs.Get("topics", "topics", 0, "id")
		prs.Get("topics", "topics", 29, "id")
	}
}

func BenchmarkFastjsonGetLarge(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var p fastjson.Parser
		prs, _ := p.ParseBytes(LargeFixture)
		prs.Get("users", "0", "id")
		prs.Get("users", "31", "id")
		prs.Get("topics", "topics", "0", "id")
		prs.Get("topics", "topics", "29", "id")
	}
}

func BenchmarkJintParseGetLarge(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		prs, _ := jin.Parse(LargeFixture)
		prs.Get("users", "0", "id")
		prs.Get("users", "31", "id")
		prs.Get("topics", "topics", "0", "id")
		prs.Get("topics", "topics", "29", "id")
	}
}