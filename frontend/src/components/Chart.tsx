import React, { useMemo, useRef, useEffect } from 'react';
import ReactECharts from 'echarts-for-react';
import { useMarketStore } from '../store/useMarketStore';
import { CHART_CONFIG } from '../constants';
import { formatDateTime } from '../utils/formatters';

const Chart: React.FC = () => {
  const chartRef = useRef<ReactECharts>(null);
  const { klines } = useMarketStore();

  const options = useMemo(() => {
    const dates = klines.map(k => formatDateTime(k.t));

    const candles = klines.map(k => [
      parseFloat(k.o),
      parseFloat(k.c),
      parseFloat(k.l),
      parseFloat(k.h)
    ]);

    const volumes = klines.map(k => parseFloat(k.v));

    // 根据数据量动态调整显示范围
    const dataLength = klines.length;
    let zoomStart = 0;
    if (dataLength > 200) {
      zoomStart = 80; // 数据多的时候显示后20%
    } else if (dataLength > 100) {
      zoomStart = 70; // 数据中等的时候显示后30%
    } else if (dataLength > 50) {
      zoomStart = 50; // 数据少的时候显示后50%
    } else {
      zoomStart = 0;  // 数据很少的时候显示全部
    }

    return {
      backgroundColor: CHART_CONFIG.BACKGROUND_COLOR,
      animation: true,
      animationDuration: 300,
      legend: {
        bottom: 10,
        left: 'center',
        data: ['K-Line', 'Volume'],
        textStyle: { color: '#ccc' }
      },
      tooltip: {
        trigger: 'axis',
        axisPointer: { type: 'cross' },
        backgroundColor: 'rgba(50, 50, 50, 0.9)',
        borderColor: '#444',
        textStyle: { color: '#eee' },
      },
      grid: [
        { left: '50', right: '50', height: '60%' },
        { left: '50', right: '50', top: '75%', height: '15%' }
      ],
      xAxis: [
        {
          type: 'category',
          data: dates,
          boundaryGap: true, // 改为true让K线有间距
          axisLine: { onZero: false, lineStyle: { color: '#444' } },
          splitLine: { show: false },
          min: 'dataMin',
          max: 'dataMax',
        },
        {
          type: 'category',
          gridIndex: 1,
          data: dates,
          boundaryGap: true,
          axisLine: { onZero: false, lineStyle: { color: '#444' } },
          axisTick: { show: false },
          splitLine: { show: false },
          axisLabel: { show: false },
          min: 'dataMin',
          max: 'dataMax'
        }
      ],
      yAxis: [
        {
          scale: true,
          splitArea: { show: false },
          axisLine: { lineStyle: { color: '#444' } },
          splitLine: { lineStyle: { color: '#333' } },
        },
        {
          scale: true,
          gridIndex: 1,
          splitNumber: 2,
          axisLabel: { show: false },
          axisLine: { show: false },
          axisTick: { show: false },
          splitLine: { show: false }
        }
      ],
      dataZoom: [
        {
          type: 'inside',
          xAxisIndex: [0, 1],
          start: zoomStart,
          end: 100
        },
        {
          show: true,
          xAxisIndex: [0, 1],
          type: 'slider',
          top: '92%',
          start: zoomStart,
          end: 100,
          textStyle: { color: '#888' }
        }
      ],
      series: [
        {
          name: 'K-Line',
          type: 'candlestick',
          data: candles,
          barWidth: '60%', // 设置K线宽度
          itemStyle: {
            color: CHART_CONFIG.UP_COLOR,
            color0: CHART_CONFIG.DOWN_COLOR,
            borderColor: CHART_CONFIG.UP_COLOR,
            borderColor0: CHART_CONFIG.DOWN_COLOR
          }
        },
        {
          name: 'Volume',
          type: 'bar',
          xAxisIndex: 1,
          yAxisIndex: 1,
          data: volumes,
          barWidth: '60%',
          itemStyle: {
            color: (params: { dataIndex: number }) => {
              const idx = params.dataIndex;
              if (idx === 0) return CHART_CONFIG.UP_COLOR;
              return candles[idx][1] >= candles[idx][0] ? CHART_CONFIG.UP_COLOR : CHART_CONFIG.DOWN_COLOR;
            }
          }
        }
      ]
    };
  }, [klines]);

  useEffect(() => {
    if (chartRef.current) {
      const chartInstance = chartRef.current.getEchartsInstance();
      chartInstance.resize();
    }
  }, []);

  return (
    <div className="w-full h-[600px]">
      <ReactECharts
        ref={chartRef}
        option={options}
        style={{ height: '100%', width: '100%' }}
        theme="dark"
        notMerge={true}
        lazyUpdate={true}
      />
    </div>
  );
};

export default Chart;
