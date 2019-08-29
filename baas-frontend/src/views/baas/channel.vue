<template>
  <div class="mixin-components-container">
    <sticky class-name="sub-navbar default">
      <el-button v-loading="loading" type="success" @click="exportChain">
        {{ $t('button.export') }}
      </el-button>
    </sticky>
    <el-collapse style="padding: 10px;">
      <div v-for="c in chainPods" :key="c.name">
        <el-collapse-item>
          <template slot="title">
            <div style="padding: 10px;font-size: 15px;">
              <i v-if="c.status == 'Running'" class="header-icon el-icon-success" style="color: greenyellow;" />
              <i v-if="c.status != 'Running'" class="header-icon el-icon-error" style="color:red;" />
              {{ c.name }}
            </div>
          </template>
          <div style="padding: 10px;font-size: 15px;">
            <el-row>
              <el-col :span="8"><div>状态: {{ c.status }}</div></el-col>
              <el-col :span="8"><div>类型: {{ c.type }}</div></el-col>
              <el-col :span="8"><div>创建时间：{{ c.createTime }}</div></el-col>
            </el-row>
          </div>
          <div style="padding: 10px;font-size: 15px;">
            <el-row>
              <el-col :span="8"><div>IP：{{ c.hostIP }}</div></el-col>
              <el-col :span="8"><div>端口：{{ c.port }}</div></el-col>
            </el-row>
          </div>
          <div style="padding: 10px;border:1px solid #DCDFE6;border-radius: 4px;color:red;text-align: center;font-size: 15px;">
            <el-row>
              <el-col :span="8"><div>CPU：{{ c.cpu }}</div></el-col>
              <el-col :span="8"><div>内存：{{ c.memory }}Mi</div></el-col>
              <el-col :span="8"><div><el-button type="primary" plain @click="handleChangeSize(c.name,c.cpu,c.memory)">{{ $t('button.changeSize') }}</el-button></div></el-col>
            </el-row>
          </div>
        </el-collapse-item>
      </div>
    </el-collapse>
    <el-row :gutter="20" class="channelbody">
      <!-- <div class="app-container documentation-container">
        <a class="document-btn" href="#">区块链</a>
        <a class="document-btn" href="#">高度</a>
        <a class="document-btn" href="#">通道数量</a>
        <a class="document-btn" href="#">链码数量</a>
      </div> -->
      <div v-for="o in channels" :key="o.id">
        <el-col :span="6">
          <el-card style="margin: 6px;">
            <el-container style="height:180px;">
              <el-aside style="width: 100px;background: white;"><el-button type="primary" icon="el-icon-news" circle /></el-aside>
              <el-main style="margin: -12px -30px;">
                <div>
                  <span style="font-size: 30px;">{{ o.channelName }}</span>
                  <div class="bottom clearfix">
                    <time class="time">{{ o.created | parseTime('{y}-{m}-{d} {h}:{i}') }}</time>
                    <div style="margin-top: 20px;">
                      <el-tag v-for="tag in o.orgs.split(',')" :key="tag" style="margin-right: 5px;" :disable-transitions="false"> {{ tag }} </el-tag>
                    </div>
                  </div>
                  <router-link :to="'/baas/chaincode/'+o.id">
                    <el-button type="text" class="button">{{ $t('button.see') }}</el-button>
                  </router-link>
                </div>
              </el-main>
            </el-container>
          </el-card>
        </el-col>
      </div>
      <el-col :span="6">
        <el-card class="box-card card">
          <div v-show="dialogFormVisible" style="height:180px;">
            <el-form ref="dataForm" :model="channel" :rules="channelRules">
              <el-form-item prop="channelName">
                <md-input v-model="channel.channelName" name="channelName" placeholder="输入通道名">
                  {{ $t('channel.channelName') }}
                </md-input>
              </el-form-item>
              <el-form-item prop="orgs">
                <el-drag-select v-model="orgs" style="width:100%;" multiple placeholder="请选择组织">
                  <el-option v-for="item in options" :key="item.value" :label="item.label" :value="item.value" />
                </el-drag-select>
              </el-form-item>
            </el-form>
            <div class="bottom clearfix">
              <el-button class="button" type="text" @click="createData()">{{ $t('button.add') }}</el-button>
            </div>
          </div>
          <div v-show="!dialogFormVisible" class="passbutton">
            <a @click="handleCreate()"><i class="el-icon-circle-plus-outline" /></a>
          </div>
        </el-card>
      </el-col>

    </el-row>
    <el-dialog :title="node" :visible.sync="dialogFormVisible2">
      <el-form ref="dataForm2" :model="channel" label-position="left" label-width="70px" style="width: 500px; margin-left:50px;">
        <el-form-item :label="$t('resource.cpu')">
          <el-slider v-model="cpu" :step="0.25" :max="maxCpu" show-stops show-input />
        </el-form-item>
        <el-form-item :label="$t('resource.memory')">
          <el-slider v-model="memory" :step="256" :max="maxMemory" show-stops show-input />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible2 = false">
          {{ $t('button.cancel') }}
        </el-button>
        <el-button type="primary" @click="changeSize()">
          {{ $t('button.confirm') }}
        </el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import MdInput from '@/components/MDinput'
import Sticky from '@/components/Sticky' // 粘性header组件
// import Mallki from '@/components/TextHoverEffect/Mallki'
import ElDragSelect from '@/components/DragSelect' // base on element-ui
import { fetch, download, podsQuery, changeSize } from '@/api/chain'
import { fetchAllList, add } from '@/api/channel'
import { mapGetters } from 'vuex'
import { parseTime } from '@/utils'

let vm = {}
export default {
  name: 'Channel',
  components: {
    MdInput,
    ElDragSelect,
    Sticky
    // Mallki
  },
  data() {
    vm = this
    // 定时刷新
    var intervalId = setInterval(function() {
      var path = sessionStorage.getItem('ROUTE_PATH')
      if (path.indexOf('/baas/channel') !== 0) {
        clearInterval(intervalId)
      }
      podsQuery(vm.chainId).then(response => {
        if (response.code === 0) {
          vm.chainPods = response.data
        }
      })
    }, 10 * 1000)

    return {
      chainId: 0,
      channel: {
        chainId: 0,
        orgs: '',
        channelName: '',
        userAccount: undefined
      },
      chain: {
        id: 0,
        userAccount: ''
      },
      memory: 0,
      cpu: 0,
      maxMemory: 2048,
      maxCpu: 2,
      node: '',
      chainPods: [],
      channelRules: {
        channelName: [{ required: true, message: 'channelName is required', trigger: 'blur' }],
        orgs: [{ required: true, message: 'orgs is required', trigger: 'blur' }]
      },
      options: [],
      orgs: [],
      channels: [],
      loading: false,
      dialogFormVisible: false,
      dialogFormVisible2: false,
      fullscreenLoading: false
    }
  },
  computed: {
    ...mapGetters([
      'account'
    ])
  },
  created() {
    const id = this.$route.params && this.$route.params.id
    this.chainId = Number(id)
    this.getChain()
    this.getAllChannel()
  },
  methods: {
    getChain(id) {
      this.chain.id = this.chainId
      this.chain.userAccount = this.account
      fetch(this.chain).then(response => {
        if (response.code === 0) {
          const orgs = response.data.peersOrgs.split(',')
          orgs.forEach(v => {
            var o = { value: v, label: v }
            this.options.push(o)
          })
        }
      })

      podsQuery(this.chainId).then(response => {
        if (response.code === 0) {
          this.chainPods = response.data
        }
      })
    },
    getAllChannel() {
      fetchAllList(this.chainId).then(response => {
        if (response.code === 0) {
          this.channels = response.data
        }
      })
    },
    exportChain() {
      this.loading = true
      if (download(this.chainId)) {
        this.loading = false
      }
    },
    handleCreate() {
      this.reset()
      this.dialogFormVisible = true
      this.$nextTick(() => {
        this.$refs['dataForm'].clearValidate()
      })
    },
    createData() {
      const loading = this.$loading({
        lock: true,
        target: '.channelbody'
      })

      this.dialogFormVisible = false
      this.channel.userAccount = this.account
      this.channel.chainId = this.chainId
      this.channel.orgs = this.orgs.join(',')
      this.$refs['dataForm'].validate((valid) => {
        if (valid) {
          add(this.channel).then(() => {
            this.$notify({
              title: '成功',
              message: '创建成功',
              type: 'success',
              duration: 2000
            })
            this.getAllChannel()
            loading.close()
          })
        }
      })
    },
    handleChangeSize(id, cpu, memory) {
      this.dialogFormVisible2 = true
      this.node = id
      this.cpu = Number(cpu)
      this.memory = Number(memory)
      this.$nextTick(() => {
        this.$refs['dataForm2'].clearValidate()
      })
    },
    changeSize() {
      var data = {
        node: this.node,
        cpu: this.cpu,
        memory: this.memory
      }
      changeSize(data).then(() => {
        this.$notify({
          title: '成功',
          message: '改变规格成功',
          type: 'success',
          duration: 2000
        })
        this.dialogFormVisible2 = false
      })
    },
    reset() {
      this.orgs = []
      this.channel = {
        chainId: 0,
        orgs: '',
        channelName: '',
        userAccount: undefined
      }
    },
    formatJson(filterVal, jsonData) {
      return jsonData.map(v => filterVal.map(j => {
        if (j === 'timestamp') {
          return parseTime(v[j])
        } else {
          return v[j]
        }
      }))
    }
  }
}
</script>

<style lang="scss" scoped>
.documentation-container {
  margin-bottom: 50px;
  .document-btn {
    float: left;
    margin-left: 100px;
    display: block;
    cursor: pointer;
    background: #409eff;
    color: white;
    height: 60px;
    width: 16%;
    line-height: 60px;
    font-size: 20px;
    text-align: center;
  }
}
</style>

<style scoped>

.mixin-components-container {
  background-color: #f0f2f5;
  padding: 30px;
  min-height: calc(100vh - 84px);
}
.component-item{
  min-height: 200px;
}

.bottom {
  margin-top: 13px;
  line-height: 12px;
}

.button {
  padding: 15px;
  float: right;
}
.passbutton {
  font-size: 130px;
  text-align: center;
  margin: 15px 0px;
  color: #1890ff;
}
.clearfix:before,
.clearfix:after {
    display: table;
    content: "";
}

.clearfix:after {
    clear: both;
}

.card {
    margin: 5px;
}

.longbutton {
    width: 173px;
}

</style>
