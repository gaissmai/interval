package interval_test

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gaissmai/interval"
)

// example time period interval
type timeInterval struct {
	birth time.Time
	death time.Time
	name  string
}

// cmp func for timeInterval
func cmpTimeInterval(p, q timeInterval) (ll, rr, lr, rl int) {
	return cmpTime(p.birth, q.birth),
		cmpTime(p.death, q.death),
		cmpTime(p.birth, q.death),
		cmpTime(p.death, q.birth)
}

// little helper
func makeTimeInterval(i, j int, s string) timeInterval {
	t1, _ := time.Parse("2006", strconv.Itoa(i))
	t2, _ := time.Parse("2006", strconv.Itoa(j))
	return timeInterval{birth: t1, death: t2, name: s}
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

// example data
var physicists = []timeInterval{
	makeTimeInterval(1473, 1543, "Kopernikus"),
	makeTimeInterval(1544, 1603, "Gilbert"),
	makeTimeInterval(1564, 1642, "Galilei"),
	makeTimeInterval(1571, 1630, "Kepler"),
	makeTimeInterval(1623, 1662, "Pascal"),
	makeTimeInterval(1629, 1695, "Huygens"),
	makeTimeInterval(1643, 1727, "Newton"),
	makeTimeInterval(1700, 1782, "Bernoulli"),
	makeTimeInterval(1777, 1855, "Gauss"),
	makeTimeInterval(1707, 1783, "Euler"),
	makeTimeInterval(1731, 1810, "Cavendish"),
	makeTimeInterval(1736, 1813, "Lagrange"),
	makeTimeInterval(1736, 1806, "Coulomb"),
	makeTimeInterval(1745, 1827, "Volta"),
	makeTimeInterval(1749, 1827, "Laplace"),
	makeTimeInterval(1768, 1830, "Fourier"),
	makeTimeInterval(1773, 1829, "Young"),
	makeTimeInterval(1775, 1836, "Ampère"),
	makeTimeInterval(1788, 1827, "Fresnel"),
	makeTimeInterval(1791, 1867, "Faraday"),
	makeTimeInterval(1796, 1832, "Carnot"),
	makeTimeInterval(1805, 1865, "Hamilton"),
	makeTimeInterval(1818, 1889, "Joule"),
	makeTimeInterval(1821, 1894, "Helholtz"),
	makeTimeInterval(1822, 1888, "Clausius"),
	makeTimeInterval(1824, 1887, "Kirchhoff"),
	makeTimeInterval(1824, 1907, "Kelvin"),
	makeTimeInterval(1831, 1879, "Maxwell"),
}

// String, implements fmt.Stringer for nice formattting
func (p timeInterval) String() string {
	return fmt.Sprintf("%s...%s (%s)", p.birth.Format("2006"), p.death.Format("2006"), p.name)
}

func ExampleTree_Precedes_time() {
	tree := interval.NewTree(cmpTimeInterval, physicists...)
	tree.Fprint(os.Stdout)

	precedes := tree.Precedes(makeTimeInterval(1643, 1727, "Newton"))
	tree = interval.NewTree(cmpTimeInterval, precedes...)

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
	tree := interval.NewTree(cmpTimeInterval, physicists...)
	tree.Fprint(os.Stdout)

	precededBy := tree.PrecededBy(makeTimeInterval(1643, 1727, "Newton"))
	tree = interval.NewTree(cmpTimeInterval, precededBy...)

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
