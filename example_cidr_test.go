package interval_test

import (
	"net/netip"
	"os"

	"github.com/gaissmai/extnetip"
	"github.com/gaissmai/interval"
)

// little helper
func mustParse(s string) MyCIDR {
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
	mustParse("10.0.0.0/8"),
	mustParse("10.32.12.1/32"),
	mustParse("10.32.16.0/22"),
	mustParse("10.32.20.0/22"),
	mustParse("10.0.0.1/32"),
	mustParse("10.0.16.0/20"),
	mustParse("10.0.32.0/20"),
	mustParse("10.0.32.1/32"),
	mustParse("10.0.48.0/20"),
	mustParse("10.1.0.0/16"),
	mustParse("10.2.0.0/15"),
	mustParse("10.4.0.0/14"),
	mustParse("10.8.0.0/13"),
	mustParse("10.16.0.0/12"),
	mustParse("10.32.0.0/11"),
	mustParse("10.32.8.0/22"),
	mustParse("10.32.12.0/22"),
	mustParse("10.64.16.0/20"),
	mustParse("10.64.32.0/19"),
	mustParse("10.64.64.0/18"),
	mustParse("10.64.128.0/17"),
	mustParse("10.65.0.0/16"),
	mustParse("10.66.0.0/15"),
	mustParse("10.64.0.0/11"),
	mustParse("10.64.0.1/32"),
	mustParse("10.64.4.0/22"),
	mustParse("fdcd:aa59:8b00::/41"),
	mustParse("fdcd:aa59:8b80::/42"),
	mustParse("fdcd:aa59:8bc0::/45"),
	mustParse("10.64.8.0/23"),
	mustParse("10.64.10.0/23"),
	mustParse("10.64.12.0/22"),
	mustParse("10.0.0.0/9"),
	mustParse("10.80.0.0/12"),
	mustParse("fdcd:aa59:8800::/39"),
	mustParse("fdcd:aa59:8a00::/40"),
	mustParse("fdcd:aa59:8bcc::/47"),
	mustParse("fdcd:aa59:8bc8::/46"),
	mustParse("fdcd:aa59:8bce::/48"),
	mustParse("fdcd:aa59:8bce::/56"),
	mustParse("fc00::/7"),
	mustParse("fdcd:aa59::/32"),
	mustParse("10.68.0.0/14"),
	mustParse("10.32.24.0/21"),
	mustParse("10.34.0.0/15"),
	mustParse("10.36.0.0/14"),
	mustParse("10.40.0.0/13"),
	mustParse("fdcd:aa59:8000::/37"),
	mustParse("10.48.0.0/12"),
	mustParse("fdcd:aa59:8bce:8::/61"),
	mustParse("fdcd:aa59:8bce:10::/60"),
	mustParse("fdcd:aa59:8bce:20::/59"),
	mustParse("fdcd:aa59:8bce:40::/58"),
	mustParse("fdcd:aa59:8bce:80::/57"),
	mustParse("fdcd:aa59:8bce:4::/62"),
	mustParse("fdcd:aa59:8bce:800::/53"),
	mustParse("fdcd:aa59:8bce:1000::/52"),
	mustParse("fdcd:aa59:8bce:2000::/51"),
	mustParse("fdcd:aa59:8bce:4000::/50"),
	mustParse("fdcd:aa59:8bce:8000::/49"),
	mustParse("fdcd:aa59:8bce:100::/56"),
	mustParse("fdcd:aa59:8bce:200::/55"),
	mustParse("fdcd:aa59:8bce:400::/54"),
	mustParse("fdcd:aa59:8bce::/64"),
	mustParse("fdcd:aa59:8bce:1::/64"),
	mustParse("fdcd:aa59:8bce:2::/64"),
	mustParse("10.0.0.0/11"),
	mustParse("fdcd:aa59:8bce:3::/64"),
	mustParse("10.72.0.0/13"),
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
