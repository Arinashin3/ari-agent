package main

import "github.com/Arinashin3/ari-agent/cmd/unisphere_exporter/collector"

//TIP <p>To run your code, right-click the code and select <b>RunMeter</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>RunMeter</b> menu item from here.</p>

//var (
//	cfgFile = kingpin.Flag("config.file", "Path to config file.").Short('c').Default("config.yml").String()
//	//	listen  = kingpin.Flag("listen", "Address to listen on").Short('l').Default(":9748").String()
//)

func main() {

	go collector.Run()

	//go collector.SystemProvider()

	select {}
}
