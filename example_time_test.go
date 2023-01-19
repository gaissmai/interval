package interval_test

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gaissmai/interval"
)

// little helper
func mkTval(i, j int, s string) Tval {
	t1, _ := time.Parse("2006", strconv.Itoa(i))
	t2, _ := time.Parse("2006", strconv.Itoa(j))
	return Tval{birth: t1, death: t2, name: s}
}

// little helper
func cmpTime(a, b time.Time) int {
	if a.Before(b) {
		return -1
	}
	if a.After(b) {
		return 1
	}
	return 0
}

// example time period interval
type Tval struct {
	birth time.Time
	death time.Time
	name  string
}

// String, implements fmt.Stringer for nice formattting, not required for interval.Interface
func (p Tval) String() string {
	return fmt.Sprintf("%s...%s (%s)", p.birth.Format("2006"), p.death.Format("2006"), p.name)
}

// Compare implements the interval.Interface
func (p Tval) Compare(q Tval) (ll, rr, lr, rl int) {
	return cmpTime(p.birth, q.birth),
		cmpTime(p.death, q.death),
		cmpTime(p.birth, q.death),
		cmpTime(p.death, q.birth)
}

// example data
var physicists = []Tval{
	mkTval(1473, 1543, "Kopernikus"),
	mkTval(1544, 1603, "Gilbert"),
	mkTval(1564, 1642, "Galilei"),
	mkTval(1571, 1630, "Kepler"),
	mkTval(1623, 1662, "Pascal"),
	mkTval(1629, 1695, "Huygens"),
	mkTval(1643, 1727, "Newton"),
	mkTval(1700, 1782, "Bernoulli"),
	mkTval(1777, 1855, "Gauss"),
	mkTval(1707, 1783, "Euler"),
	mkTval(1731, 1810, "Cavendish"),
	mkTval(1736, 1813, "Lagrange"),
	mkTval(1736, 1806, "Coulomb"),
	mkTval(1745, 1827, "Volta"),
	mkTval(1749, 1827, "Laplace"),
	mkTval(1768, 1830, "Fourier"),
	mkTval(1773, 1829, "Young"),
	mkTval(1775, 1836, "Ampère"),
	mkTval(1788, 1827, "Fresnel"),
	mkTval(1791, 1867, "Faraday"),
	mkTval(1796, 1832, "Carnot"),
	mkTval(1805, 1865, "Hamilton"),
	mkTval(1818, 1889, "Joule"),
	mkTval(1821, 1894, "Helholtz"),
	mkTval(1822, 1888, "Clausius"),
	mkTval(1824, 1887, "Kirchhoff"),
	mkTval(1824, 1907, "Kelvin"),
	mkTval(1831, 1879, "Maxwell"),
}

func ExampleTree_Precedes_time() {
	tree := interval.NewTree(physicists...)
	tree.Fprint(os.Stdout)

	precedes := tree.Precedes(mkTval(1643, 1727, "Newton"))
	tree = interval.NewTree(precedes...)

	fmt.Println("\nPrecedes Newton:")
	tree.Fprint(os.Stdout)

	// Output:
	// ▼
	// ├─ 1473...1543 (Kopernikus)
	// ├─ 1544...1603 (Gilbert)
	// ├─ 1564...1642 (Galilei)
	// │  └─ 1571...1630 (Kepler)
	// ├─ 1623...1662 (Pascal)
	// ├─ 1629...1695 (Huygens)
	// ├─ 1643...1727 (Newton)
	// ├─ 1700...1782 (Bernoulli)
	// ├─ 1707...1783 (Euler)
	// ├─ 1731...1810 (Cavendish)
	// ├─ 1736...1813 (Lagrange)
	// │  └─ 1736...1806 (Coulomb)
	// ├─ 1745...1827 (Volta)
	// │  └─ 1749...1827 (Laplace)
	// ├─ 1768...1830 (Fourier)
	// │  └─ 1773...1829 (Young)
	// ├─ 1775...1836 (Ampère)
	// ├─ 1777...1855 (Gauss)
	// │  └─ 1788...1827 (Fresnel)
	// ├─ 1791...1867 (Faraday)
	// │  ├─ 1796...1832 (Carnot)
	// │  └─ 1805...1865 (Hamilton)
	// ├─ 1818...1889 (Joule)
	// ├─ 1821...1894 (Helholtz)
	// │  └─ 1822...1888 (Clausius)
	// └─ 1824...1907 (Kelvin)
	//    └─ 1824...1887 (Kirchhoff)
	//       └─ 1831...1879 (Maxwell)
	//
	// Precedes Newton:
	// ▼
	// ├─ 1473...1543 (Kopernikus)
	// ├─ 1544...1603 (Gilbert)
	// └─ 1564...1642 (Galilei)
	//    └─ 1571...1630 (Kepler)
}

func ExampleTree_PrecededBy_time() {
	tree := interval.NewTree(physicists...)
	tree.Fprint(os.Stdout)

	precededBy := tree.PrecededBy(mkTval(1643, 1727, "Newton"))
	tree = interval.NewTree(precededBy...)

	fmt.Println("\nPrecededBy Newton:")
	tree.Fprint(os.Stdout)

	// Output:
	// ▼
	// ├─ 1473...1543 (Kopernikus)
	// ├─ 1544...1603 (Gilbert)
	// ├─ 1564...1642 (Galilei)
	// │  └─ 1571...1630 (Kepler)
	// ├─ 1623...1662 (Pascal)
	// ├─ 1629...1695 (Huygens)
	// ├─ 1643...1727 (Newton)
	// ├─ 1700...1782 (Bernoulli)
	// ├─ 1707...1783 (Euler)
	// ├─ 1731...1810 (Cavendish)
	// ├─ 1736...1813 (Lagrange)
	// │  └─ 1736...1806 (Coulomb)
	// ├─ 1745...1827 (Volta)
	// │  └─ 1749...1827 (Laplace)
	// ├─ 1768...1830 (Fourier)
	// │  └─ 1773...1829 (Young)
	// ├─ 1775...1836 (Ampère)
	// ├─ 1777...1855 (Gauss)
	// │  └─ 1788...1827 (Fresnel)
	// ├─ 1791...1867 (Faraday)
	// │  ├─ 1796...1832 (Carnot)
	// │  └─ 1805...1865 (Hamilton)
	// ├─ 1818...1889 (Joule)
	// ├─ 1821...1894 (Helholtz)
	// │  └─ 1822...1888 (Clausius)
	// └─ 1824...1907 (Kelvin)
	//    └─ 1824...1887 (Kirchhoff)
	//       └─ 1831...1879 (Maxwell)
	//
	// PrecededBy Newton:
	// ▼
	// ├─ 1731...1810 (Cavendish)
	// ├─ 1736...1813 (Lagrange)
	// │  └─ 1736...1806 (Coulomb)
	// ├─ 1745...1827 (Volta)
	// │  └─ 1749...1827 (Laplace)
	// ├─ 1768...1830 (Fourier)
	// │  └─ 1773...1829 (Young)
	// ├─ 1775...1836 (Ampère)
	// ├─ 1777...1855 (Gauss)
	// │  └─ 1788...1827 (Fresnel)
	// ├─ 1791...1867 (Faraday)
	// │  ├─ 1796...1832 (Carnot)
	// │  └─ 1805...1865 (Hamilton)
	// ├─ 1818...1889 (Joule)
	// ├─ 1821...1894 (Helholtz)
	// │  └─ 1822...1888 (Clausius)
	// └─ 1824...1907 (Kelvin)
	//    └─ 1824...1887 (Kirchhoff)
	//       └─ 1831...1879 (Maxwell)
}
