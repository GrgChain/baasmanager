<template>
  <div class="dashboard-editor-container">
    <panel-group @handleSetLineChartData="handleSetLineChartData" />

    <el-row style="background:#fff;padding:16px 16px 0;margin-bottom:32px;">
      <line-chart :chart-data="lineChartData" />
    </el-row>

    <el-row :gutter="32">
      <el-col :xs="24" :sm="24" :lg="8">
        <div class="chart-wrapper">
          <pie-chart />
        </div>
      </el-col>
      <!-- <el-col :xs="24" :sm="24" :lg="8">
        <div class="chart-wrapper">
          <raddar-chart />
        </div>
      </el-col>
      <el-col :xs="24" :sm="24" :lg="8">
        <div class="chart-wrapper">
          <bar-chart />
        </div>
      </el-col> -->
    </el-row>
  </div>
</template>

<script>
import PanelGroup from './components/PanelGroup'
import LineChart from './components/LineChart'
// import RaddarChart from './components/RaddarChart'
import PieChart from './components/PieChart'
// import BarChart from './components/BarChart'

import { getSevenDays } from '@/api/dash'
import { mapGetters } from 'vuex'
import { onlySevenDays, sevenDays } from '@/utils'

const lineChartData = {
  users: {
    actualData: [0, 0, 0, 0, 0, 0, 0]
  },
  chains: {
    actualData: [0, 0, 0, 0, 0, 0, 0]
  },
  channels: {
    actualData: [0, 0, 0, 0, 0, 0, 0]
  },
  chaincodes: {
    actualData: [0, 0, 0, 0, 0, 0, 0]
  }
}
function parseList(data) {
  var days = sevenDays()[1]
  var users = data['users']
  var chains = data['chains']
  var channels = data['channels']
  var chaincodes = data['chaincodes']

  var i = 0; var j = 0; var counts = []; var detail = {}

  if (users != null) {
    counts = [0, 0, 0, 0, 0, 0, 0]
    for (j = 0; j < users.length; j++) {
      detail = users[j]
      for (i = 0; i < 7; i++) {
        if (detail.days === days[i]) {
          counts[i] = detail.counts
        }
      }
    }
    lineChartData.users.actualData = counts
  }

  if (chaincodes != null) {
    counts = [0, 0, 0, 0, 0, 0, 0]
    for (j = 0; j < chaincodes.length; j++) {
      detail = chaincodes[j]
      for (i = 0; i < 7; i++) {
        if (detail.days === days[i]) {
          counts[i] = detail.counts
        }
      }
    }
    lineChartData.chaincodes.actualData = counts
  }

  if (channels != null) {
    counts = [0, 0, 0, 0, 0, 0, 0]
    for (j = 0; j < channels.length; j++) {
      detail = channels[j]
      for (i = 0; i < 7; i++) {
        if (detail.days === days[i]) {
          counts[i] = detail.counts
        }
      }
    }
    lineChartData.channels.actualData = counts
  }

  if (chains != null) {
    counts = [0, 0, 0, 0, 0, 0, 0]
    for (j = 0; j < chains.length; j++) {
      detail = chains[j]
      for (i = 0; i < 7; i++) {
        if (detail.days === days[i]) {
          counts[i] = detail.counts
        }
      }
    }
    lineChartData.chains.actualData = counts
  }
}

export default {
  name: 'DashboardAdmin',
  components: {
    PanelGroup,
    LineChart,
    // RaddarChart,
    // BarChart,
    PieChart
  },
  data() {
    return {
      query: {
        start: 0,
        end: 0,
        userAccount: undefined
      },
      lineChartData: lineChartData.users
    }
  },
  computed: {
    ...mapGetters([
      'account'
    ])
  },
  created() {
    var days = onlySevenDays()
    this.query.start = days.start
    this.query.end = days.end
    this.getDays()
  },
  methods: {
    handleSetLineChartData(type) {
      this.lineChartData = lineChartData[type]
    },
    getDays() {
      this.query.userAccount = this.account
      getSevenDays(this.query).then(response => {
        // Just to simulate the time of the request
        parseList(response.data)
      })
    }
  }
}
</script>

<style lang="scss" scoped>
.dashboard-editor-container {
  padding: 32px;
  background-color: rgb(240, 242, 245);
  position: relative;

  .github-corner {
    position: absolute;
    top: 0px;
    border: 0;
    right: 0;
  }

  .chart-wrapper {
    background: #fff;
    padding: 16px 16px 0;
    margin-bottom: 32px;
  }
}
</style>
