(function () {
  "use strict";

  var el = document.getElementById("chart-modules");
  if (!el || typeof echarts === "undefined") return;

  var chart = echarts.init(el);

  var categories = [
    "可源码重编",
    "官方预编译",
    "Windows 专属",
    "闭源 darwin 专属",
    "darwin 子应用"
  ];
  var values = [7, 2, 4, 5, 1];
  var detail = [
    "sqlite3 · spdlog · node-pty · watchdog · keymap · is-elevated · kerberos",
    "@parcel/watcher · rg 搜索二进制",
    "deviceid · registry-js · foreground-love · policy-watcher",
    "ipc · net · perf-sdk · network-client · macos-native",
    "mac-computer-use（MCP 子应用）"
  ];
  var colors = ["#0b66e3", "#3b82c4", "#8fa3b8", "#d95a1e", "#b8860b"];

  chart.setOption({
    baseOption: {
      color: colors,
      tooltip: {
        trigger: "item",
        confine: true,
        backgroundColor: "#ffffff",
        borderColor: "#dfe4ec",
        textStyle: { color: "#1c2733", fontSize: 12 },
        formatter: function (p) {
          return (
            "<b>" + p.name + "</b> — " + p.value + " 个<br/>" +
            '<span style="color:#5c6b7a">' + detail[p.dataIndex] + "</span>"
          );
        }
      },
      legend: {
        bottom: 0,
        icon: "roundRect",
        itemWidth: 12,
        itemHeight: 8,
        textStyle: { color: "#5c6b7a", fontSize: 12 }
      },
      series: [
        {
          name: "原生二进制分布",
          type: "pie",
          radius: ["38%", "62%"],
          center: ["50%", "44%"],
          itemStyle: { borderColor: "#f6f7f9", borderWidth: 3 },
          label: {
            color: "#1c2733",
            fontSize: 12,
            formatter: "{b}\n{c} 个 ({d}%)"
          },
          labelLine: { length: 12, length2: 8, lineStyle: { color: "#dfe4ec" } },
          emphasis: {
            scale: true,
            scaleSize: 6,
            itemStyle: { shadowBlur: 12, shadowColor: "rgba(28,39,51,0.18)" }
          },
          data: categories.map(function (name, i) {
            return { name: name, value: values[i] };
          })
        }
      ]
    }
  });

  window.addEventListener("resize", function () {
    chart.resize();
  });
})();
