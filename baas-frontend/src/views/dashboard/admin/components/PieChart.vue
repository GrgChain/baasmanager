<template>
  <div :class="className" :style="{height:height,width:width}" />
</template>

<script>
import echarts from 'echarts'
require('echarts/theme/macarons') // echarts theme
import { debounce } from '@/utils'
import { consensusTotal } from '@/api/dash'
import { mapGetters } from 'vuex'

export default {
  props: {
    className: {
      type: String,
      default: 'chart'
    },
    width: {
      type: String,
      default: '100%'
    },
    height: {
      type: String,
      default: '300px'
    }
  },
  data() {
    return {
      query: {
        userAccount: undefined
      },
      totals: [
        { value: 0, name: 'Solo' },
        { value: 0, name: 'Kafka' },
        { value: 0, name: 'Etcdraft' }
      ],
      chart: null
    }
  },
  computed: {
    ...mapGetters([
      'account'
    ])
  },
  mounted() {
    this.query.userAccount = this.account
    consensusTotal(this.query).then(response => {
      // Just to simulate the time of the request
      response.data.forEach(v => {
        var tkv = {}
        if (v['consensus'] === 'solo') {
          tkv['name'] = 'Solo'
          tkv['value'] = Number(v['value'])
          this.totals[0] = tkv
        } else if (v['consensus'] === 'kafka') {
          tkv['name'] = 'Kafka'
          tkv['value'] = Number(v['value'])
          this.totals[1] = tkv
        } else if (v['consensus'] === 'etcdraft') {
          tkv['name'] = 'Etcdraft'
          tkv['value'] = Number(v['value'])
          this.totals[2] = tkv
        }
      })
      console.log(this.totals)
      this.initChart()
      this.__resizeHandler = debounce(() => {
        if (this.chart) {
          this.chart.resize()
        }
      }, 100)
      window.addEventListener('resize', this.__resizeHandler)
    })
  },
  beforeDestroy() {
    if (!this.chart) {
      return
    }
    window.removeEventListener('resize', this.__resizeHandler)
    this.chart.dispose()
    this.chart = null
  },
  created() {
  },
  methods: {

    initChart() {
      this.chart = echarts.init(this.$el, 'macarons')
      this.chart.setOption({
        tooltip: {
          trigger: 'item',
          formatter: '{a} <br/>{b} : {c} ({d}%)'
        },
        legend: {
          left: 'center',
          bottom: '10',
          data: ['Solo', 'Kafka', 'Etcdraft']
        },
        calculable: true,
        series: [
          {
            name: 'NETWORK CONSENSUS',
            type: 'pie',
            roseType: 'radius',
            radius: [15, 95],
            center: ['50%', '38%'],
            data: this.totals,
            animationEasing: 'cubicInOut',
            animationDuration: 2600
          }
        ]
      })
    }
  }
}
</script>
