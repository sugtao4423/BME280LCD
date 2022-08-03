package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/quhar/bme280"
	"golang.org/x/exp/io/i2c"

	"github.com/sugtao4423/BME280LCD/acm1602ni"
)

var (
	i2cDev      string
	bme280Addr  int
	acm1602Addr int
)

func main() {
	i2cDev = *flag.String("i2c", "/dev/i2c-1", "i2c device")
	bme280Addr = *flag.Int("bme280", 0x76, "bme280 address")
	acm1602Addr = *flag.Int("acm1602", 0x50, "acm1602ni address")
	flag.Parse()

	t, p, h := getBme280()

	time := time.Now()
	l1 := fmt.Sprintf("%.2fßC  %.2f%%", t, h)
	l2 := fmt.Sprintf("%.2fhPa %02d:%02d", p, time.Hour(), time.Minute())

	showLcd(l1, l2)
}

func getBme280() (temp, press, hum float64) {
	d, err := i2c.Open(&i2c.Devfs{Dev: i2cDev}, bme280Addr)
	if err != nil {
		panic(err)
	}
	defer d.Close()

	b := bme280.New(d)
	err = b.Init()
	if err != nil {
		panic(err)
	}
	t, p, h, err := b.EnvData()
	if err != nil {
		panic(err)
	}
	return t, p, h
}

func showLcd(line1, line2 string) {
	d, err := i2c.Open(&i2c.Devfs{Dev: i2cDev}, acm1602Addr)
	if err != nil {
		panic(err)
	}
	defer d.Close()

	a := acm1602ni.New(d)
	err = a.Init()
	if err != nil {
		panic(err)
	}

	err = a.Show(line1, line2)
	if err != nil {
		panic(err)
	}
}
