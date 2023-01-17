package interval_test

import (
	"net/netip"
	"os"

	"github.com/gaissmai/extnetip"
	"github.com/gaissmai/interval"
)

// little helper
func mustParsePrefix(s string) MyCIDR {
	pfx, err := netip.ParsePrefix(s)
	if err != nil {
		panic(err)
	}
	return MyCIDR{pfx}
}

type MyCIDR struct{ netip.Prefix }

// implement the interval.Interface for MyCIDR
func (p MyCIDR) Compare(q MyCIDR) (ll, rr, lr, rl int) {
	pl, pr := extnetip.Range(p.Prefix)
	ql, qr := extnetip.Range(q.Prefix)

	return pl.Compare(ql),
		pr.Compare(qr),
		pl.Compare(qr),
		pr.Compare(ql)
}

// example data
var cidrs = []MyCIDR{
	mustParsePrefix("10.0.0.0/8"),
	mustParsePrefix("10.32.12.1/32"),
	mustParsePrefix("10.32.16.0/22"),
	mustParsePrefix("10.32.20.0/22"),
	mustParsePrefix("10.0.0.1/32"),
	mustParsePrefix("10.0.16.0/20"),
	mustParsePrefix("10.0.32.0/20"),
	mustParsePrefix("10.0.32.1/32"),
	mustParsePrefix("10.0.48.0/20"),
	mustParsePrefix("10.1.0.0/16"),
	mustParsePrefix("10.2.0.0/15"),
	mustParsePrefix("10.4.0.0/14"),
	mustParsePrefix("10.8.0.0/13"),
	mustParsePrefix("10.16.0.0/12"),
	mustParsePrefix("10.32.0.0/11"),
	mustParsePrefix("10.32.8.0/22"),
	mustParsePrefix("10.32.12.0/22"),
	mustParsePrefix("10.64.16.0/20"),
	mustParsePrefix("10.64.32.0/19"),
	mustParsePrefix("10.64.64.0/18"),
	mustParsePrefix("10.64.128.0/17"),
	mustParsePrefix("10.65.0.0/16"),
	mustParsePrefix("10.66.0.0/15"),
	mustParsePrefix("10.64.0.0/11"),
	mustParsePrefix("10.64.0.1/32"),
	mustParsePrefix("10.64.4.0/22"),
	mustParsePrefix("fdcd:aa59:8b00::/41"),
	mustParsePrefix("fdcd:aa59:8b80::/42"),
	mustParsePrefix("fdcd:aa59:8bc0::/45"),
	mustParsePrefix("10.64.8.0/23"),
	mustParsePrefix("10.64.10.0/23"),
	mustParsePrefix("10.64.12.0/22"),
	mustParsePrefix("10.0.0.0/9"),
	mustParsePrefix("10.80.0.0/12"),
	mustParsePrefix("fdcd:aa59:8800::/39"),
	mustParsePrefix("fdcd:aa59:8a00::/40"),
	mustParsePrefix("fdcd:aa59:8bcc::/47"),
	mustParsePrefix("fdcd:aa59:8bc8::/46"),
	mustParsePrefix("fdcd:aa59:8bce::/48"),
	mustParsePrefix("fdcd:aa59:8bce::/56"),
	mustParsePrefix("fc00::/7"),
	mustParsePrefix("fdcd:aa59::/32"),
	mustParsePrefix("10.68.0.0/14"),
	mustParsePrefix("10.32.24.0/21"),
	mustParsePrefix("10.34.0.0/15"),
	mustParsePrefix("10.36.0.0/14"),
	mustParsePrefix("10.40.0.0/13"),
	mustParsePrefix("fdcd:aa59:8000::/37"),
	mustParsePrefix("10.48.0.0/12"),
	mustParsePrefix("fdcd:aa59:8bce:8::/61"),
	mustParsePrefix("fdcd:aa59:8bce:10::/60"),
	mustParsePrefix("fdcd:aa59:8bce:20::/59"),
	mustParsePrefix("fdcd:aa59:8bce:40::/58"),
	mustParsePrefix("fdcd:aa59:8bce:80::/57"),
	mustParsePrefix("fdcd:aa59:8bce:4::/62"),
	mustParsePrefix("fdcd:aa59:8bce:800::/53"),
	mustParsePrefix("fdcd:aa59:8bce:1000::/52"),
	mustParsePrefix("fdcd:aa59:8bce:2000::/51"),
	mustParsePrefix("fdcd:aa59:8bce:4000::/50"),
	mustParsePrefix("fdcd:aa59:8bce:8000::/49"),
	mustParsePrefix("fdcd:aa59:8bce:100::/56"),
	mustParsePrefix("fdcd:aa59:8bce:200::/55"),
	mustParsePrefix("fdcd:aa59:8bce:400::/54"),
	mustParsePrefix("fdcd:aa59:8bce::/64"),
	mustParsePrefix("fdcd:aa59:8bce:1::/64"),
	mustParsePrefix("fdcd:aa59:8bce:2::/64"),
	mustParsePrefix("10.0.0.0/11"),
	mustParsePrefix("fdcd:aa59:8bce:3::/64"),
	mustParsePrefix("10.72.0.0/13"),
}

func ExampleInterface_cidr() {
	tree := interval.NewTree(cidrs...)
	tree.Fprint(os.Stdout)

	// Output:
	// ▼
	// ├─ 10.0.0.0/8
	// │  └─ 10.0.0.0/9
	// │     ├─ 10.0.0.0/11
	// │     │  ├─ 10.0.0.1/32
	// │     │  ├─ 10.0.16.0/20
	// │     │  ├─ 10.0.32.0/20
	// │     │  │  └─ 10.0.32.1/32
	// │     │  ├─ 10.0.48.0/20
	// │     │  ├─ 10.1.0.0/16
	// │     │  ├─ 10.2.0.0/15
	// │     │  ├─ 10.4.0.0/14
	// │     │  ├─ 10.8.0.0/13
	// │     │  └─ 10.16.0.0/12
	// │     ├─ 10.32.0.0/11
	// │     │  ├─ 10.32.8.0/22
	// │     │  ├─ 10.32.12.0/22
	// │     │  │  └─ 10.32.12.1/32
	// │     │  ├─ 10.32.16.0/22
	// │     │  ├─ 10.32.20.0/22
	// │     │  ├─ 10.32.24.0/21
	// │     │  ├─ 10.34.0.0/15
	// │     │  ├─ 10.36.0.0/14
	// │     │  ├─ 10.40.0.0/13
	// │     │  └─ 10.48.0.0/12
	// │     └─ 10.64.0.0/11
	// │        ├─ 10.64.0.1/32
	// │        ├─ 10.64.4.0/22
	// │        ├─ 10.64.8.0/23
	// │        ├─ 10.64.10.0/23
	// │        ├─ 10.64.12.0/22
	// │        ├─ 10.64.16.0/20
	// │        ├─ 10.64.32.0/19
	// │        ├─ 10.64.64.0/18
	// │        ├─ 10.64.128.0/17
	// │        ├─ 10.65.0.0/16
	// │        ├─ 10.66.0.0/15
	// │        ├─ 10.68.0.0/14
	// │        ├─ 10.72.0.0/13
	// │        └─ 10.80.0.0/12
	// └─ fc00::/7
	//    └─ fdcd:aa59::/32
	//       ├─ fdcd:aa59:8000::/37
	//       ├─ fdcd:aa59:8800::/39
	//       ├─ fdcd:aa59:8a00::/40
	//       ├─ fdcd:aa59:8b00::/41
	//       ├─ fdcd:aa59:8b80::/42
	//       ├─ fdcd:aa59:8bc0::/45
	//       ├─ fdcd:aa59:8bc8::/46
	//       ├─ fdcd:aa59:8bcc::/47
	//       └─ fdcd:aa59:8bce::/48
	//          ├─ fdcd:aa59:8bce::/56
	//          │  ├─ fdcd:aa59:8bce::/64
	//          │  ├─ fdcd:aa59:8bce:1::/64
	//          │  ├─ fdcd:aa59:8bce:2::/64
	//          │  ├─ fdcd:aa59:8bce:3::/64
	//          │  ├─ fdcd:aa59:8bce:4::/62
	//          │  ├─ fdcd:aa59:8bce:8::/61
	//          │  ├─ fdcd:aa59:8bce:10::/60
	//          │  ├─ fdcd:aa59:8bce:20::/59
	//          │  ├─ fdcd:aa59:8bce:40::/58
	//          │  └─ fdcd:aa59:8bce:80::/57
	//          ├─ fdcd:aa59:8bce:100::/56
	//          ├─ fdcd:aa59:8bce:200::/55
	//          ├─ fdcd:aa59:8bce:400::/54
	//          ├─ fdcd:aa59:8bce:800::/53
	//          ├─ fdcd:aa59:8bce:1000::/52
	//          ├─ fdcd:aa59:8bce:2000::/51
	//          ├─ fdcd:aa59:8bce:4000::/50
	//          └─ fdcd:aa59:8bce:8000::/49
}
