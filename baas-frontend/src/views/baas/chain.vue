<template>
  <div class="app-container">
    <div class="filter-container">
      <el-input v-model="listQuery.name" :placeholder="$t('chain.name')" style="width: 200px;" class="filter-item" @keyup.enter.native="handleFilter" />
      <el-button v-waves class="filter-item" type="primary" icon="el-icon-search" @click="handleFilter">
        {{ $t('button.search') }}
      </el-button>
      <el-button class="filter-item" style="margin-left: 10px;" type="primary" icon="el-icon-plus" @click="handleCreate">
        {{ $t('button.add') }}
      </el-button>
    </div>

    <el-table
      :key="tableKey"
      v-loading="listLoading"
      :data="list"
      border
      fit
      highlight-current-row
      style="width: 100%;"
    >
      <el-table-column :label="$t('chain.id')" prop="id" align="center" width="80">
        <template slot-scope="scope">
          <span>{{ scope.row.id }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('chain.name')" min-width="150px">
        <template slot-scope="{row}">
          <span v-if="row.status==0 || row.status==1 || row.status ==3">{{ row.name }}</span>
          <router-link :to="'/baas/channel/'+row.id">
            <span v-if="row.status ==2" class="link-type">{{ row.name }}</span>
          </router-link>
        </template>
      </el-table-column>
      <el-table-column :label="$t('chain.consensus')" min-width="80px">
        <template slot-scope="{row}">
          <span>{{ row.consensus }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('chain.peersOrgs')" min-width="150px">
        <template slot-scope="{row}">
          <el-tag v-for="tag in row.peersOrgs.split(',')" :key="tag" style="margin-right: 5px;" :disable-transitions="false"> {{ tag }} </el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="$t('chain.tlsEnabled')" min-width="80px">
        <template slot-scope="{row}">
          <span>{{ row.tlsEnabled | tlsEnabledFilter }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('chain.orderCount')" min-width="80px">
        <template slot-scope="{row}">
          <span>{{ row.orderCount }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('chain.peerCount')" min-width="80px">
        <template slot-scope="{row}">
          <span>{{ row.peerCount }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('chain.description')" min-width="150px">
        <template slot-scope="{row}">
          <span>{{ row.description }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('chain.userAccount')" min-width="100px">
        <template slot-scope="{row}">
          <span>{{ row.userAccount }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('chain.created')" min-width="120px">
        <template slot-scope="{row}">
          <span>{{ row.created | parseTime('{y}-{m}-{d} {h}:{i}') }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('chain.status')" min-width="80px">
        <template slot-scope="{row}">
          <el-tag v-if="row.status==0" type="info">定义</el-tag>
          <el-tag v-if="row.status==1" type="warning">已构建</el-tag>
          <el-tag v-if="row.status==2" type="success">运行中</el-tag>
          <el-tag v-if="row.status==3" type="danger">已停止</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="$t('button.actions')" align="center" width="230" class-name="small-padding fixed-width">
        <template slot-scope="{row}">
          <el-button v-if="row.status==0" type="primary" size="mini" @click="handleUpdate(row)">
            {{ $t('button.edit') }}
          </el-button>
          <el-button v-if="row.status==0" size="mini" type="danger" @click="handleDelete(row)">
            {{ $t('button.delete') }}
          </el-button>
          <el-button v-if="row.status==0" size="mini" type="warning" @click="handleBuild(row)">
            {{ $t('button.build') }}
          </el-button>
          <el-button v-if="row.status==1 || row.status==3" size="mini" type="success" @click="handleRun(row)">
            {{ $t('button.run') }}
          </el-button>
          <el-button v-if="row.status==2" size="mini" type="info" @click="handleStop(row)">
            {{ $t('button.stop') }}
          </el-button>
          <el-button v-if="row.status==3" size="mini" type="danger" @click="handleRelease(row)">
            {{ $t('button.release') }}
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <pagination v-show="total>0" :total="total" :page.sync="listQuery.page" :limit.sync="listQuery.limit" @pagination="getList" />

    <el-dialog :title="textMap[dialogStatus]" :visible.sync="dialogFormVisible">
      <el-form ref="dataForm" :rules="rules" :model="temp" label-position="left" label-width="120px" style="width: 500px; margin-left:50px;">

        <el-form-item :label="$t('chain.name')" prop="name">
          <el-input v-model="temp.name" />
        </el-form-item>
        <el-form-item :label="$t('chain.consensus')" prop="consensus">
          <el-select v-model="temp.consensus" placeholder="请选择" @change="consensusChange">
            <el-option v-for="item in consensus" :key="item.value" :label="item.label" :value="item.value" :disabled="item.disabled" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('chain.peersOrgs')" prop="peersOrgs">
          <el-tag v-for="tag in orgTags" :key="tag" closable :disable-transitions="false" @close="handleClose(tag)">{{ tag }}</el-tag>
          <el-input v-if="inputVisible" ref="saveTagInput" v-model="inputValue" size="small" class="input-new-tag" @keyup.enter.native="handleInputConfirm" @blur="handleInputConfirm" />
          <el-button v-else class="button-new-tag" size="small" @click="showInput">+ New Org</el-button>
        </el-form-item>
        <el-form-item :label="$t('chain.tlsEnabled')" prop="tlsEnabled">
          <el-radio v-model="temp.tlsEnabled" label="true" border> {{ $t('button.true') }} </el-radio>
          <el-radio v-show="!raftShow" v-model="temp.tlsEnabled" label="false" border> {{ $t('button.false') }} </el-radio>
          <el-radio v-show="raftShow" v-model="temp.tlsEnabled" label="false" border disabled> {{ $t('button.false') }} </el-radio>
        </el-form-item>
        <el-form-item :label="$t('chain.orderCount')" prop="orderCount">
          <el-input-number v-model="temp.orderCount" :min="orderMinCount" :max="orderMaxCount" />
        </el-form-item>
        <el-form-item :label="$t('chain.peerCount')" prop="peerCount">
          <el-input-number v-model="temp.peerCount" :min="1" :max="peerMaxCount" />
        </el-form-item>
        <el-form-item :label="$t('chain.description')">
          <el-input v-model="temp.description" />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible = false">
          {{ $t('button.cancel') }}
        </el-button>
        <el-button type="primary" @click="dialogStatus==='create'?createData():updateData()">
          {{ $t('button.confirm') }}
        </el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { fetchList, add, update, del, stop, run, build, release } from '@/api/chain'
import waves from '@/directive/waves' // waves directive
import Pagination from '@/components/Pagination' // secondary package based on el-pagination
import { parseTime } from '@/utils'
import { mapGetters } from 'vuex'

export default {
  name: 'ComplexTable',
  components: { Pagination },
  directives: { waves },
  filters: {
    tlsEnabledFilter(row) {
      return row === 'true' ? '是' : '否'
    }
  },
  data() {
    return {
      tableKey: 0,
      list: null,
      total: 0,
      listLoading: true,
      listQuery: {
        page: 1,
        limit: 10,
        name: undefined,
        userAccount: undefined
      },
      consensus: [{
        value: 'solo',
        label: 'solo'
      }, {
        value: 'kafka',
        label: 'kafka'
      }, {
        value: 'etcdraft',
        label: 'etcdraft'
      }],
      orderMaxCount: 10,
      orderMinCount: 1,
      peerMaxCount: 10,
      orgTags: [],
      inputVisible: false,
      inputValue: '',
      temp: {
        consensus: undefined,
        created: 0,
        description: undefined,
        id: undefined,
        name: undefined,
        orderCount: 0,
        peerCount: 0,
        peersOrgs: undefined,
        status: 0,
        tlsEnabled: undefined,
        userAccount: undefined
      },
      raftShow: false,
      dialogFormVisible: false,
      dialogStatus: '',
      textMap: {
        update: '编辑链',
        create: '添加链'
      },
      rules: {
        name: [{ required: true, message: 'name is required', trigger: 'blur' }],
        consensus: [{ required: true, message: 'consensus is required', trigger: 'blur' }],
        orderCount: [{ required: true, message: 'orderCount is required', trigger: 'blur' }],
        peerCount: [{ required: true, message: 'peerCount is required', trigger: 'blur' }],
        peersOrgs: [{ required: true, message: 'peersOrgs is required', trigger: 'blur' }],
        tlsEnabled: [{ required: true, message: 'tlsEnabled is required', trigger: 'blur' }]
      }
    }
  },
  computed: {
    ...mapGetters([
      'account'
    ])
  },
  created() {
    this.getList()
  },
  methods: {
    consensusChange() {
      if (this.temp.consensus === 'solo') {
        this.temp.orderCount = 1
        this.orderMaxCount = 1
        this.raftShow = false
      } else if (this.temp.consensus === 'etcdraft') {
        this.temp.orderCount = 3
        this.orderMinCount = 3
        this.orderMaxCount = 5
        this.raftShow = true
      } else {
        this.orderMaxCount = 10
        this.raftShow = false
      }
    },
    handleClose(tag) {
      this.orgTags.splice(this.orgTags.indexOf(tag), 1)
      this.temp.peersOrgs = this.orgTags.join(',')
    },
    showInput() {
      this.inputVisible = true
      this.$nextTick(_ => {
        this.$refs.saveTagInput.$refs.input.focus()
      })
    },
    handleInputConfirm() {
      const inputValue = this.inputValue.toLowerCase()
      if (inputValue && this.orgTags.indexOf(inputValue) === -1) {
        this.orgTags.push(inputValue)
        this.temp.peersOrgs = this.orgTags.join(',')
      }
      this.inputVisible = false
      this.inputValue = ''
    },
    getList() {
      this.listLoading = true
      this.listQuery.userAccount = this.account
      fetchList(this.listQuery).then(response => {
        this.list = response.data
        this.total = response.total
        // Just to simulate the time of the request
        setTimeout(() => {
          this.listLoading = false
        }, 1.5 * 1000)
      })
    },
    handleFilter() {
      this.listQuery.page = 1
      this.getList()
    },
    resetTemp() {
      this.temp = {
        consensus: undefined,
        created: 0,
        description: undefined,
        id: undefined,
        name: undefined,
        orderCount: 0,
        peerCount: 0,
        peersOrgs: undefined,
        status: 0,
        tlsEnabled: undefined,
        userAccount: undefined
      }
    },
    handleCreate() {
      this.resetTemp()
      this.dialogStatus = 'create'
      this.dialogFormVisible = true
      this.orgTags = []
      this.$nextTick(() => {
        this.$refs['dataForm'].clearValidate()
      })
    },
    createData() {
      this.$refs['dataForm'].validate((valid) => {
        if (valid) {
          this.temp.userAccount = this.account
          add(this.temp).then(() => {
            this.dialogFormVisible = false
            this.$notify({
              title: '成功',
              message: '创建成功',
              type: 'success',
              duration: 2000
            })
            this.getList()
          })
        }
      })
    },
    handleUpdate(row) {
      this.temp = Object.assign({}, row) // copy obj
      this.dialogStatus = 'update'
      this.dialogFormVisible = true
      this.orgTags = this.temp.peersOrgs.split(',')
      this.$nextTick(() => {
        this.$refs['dataForm'].clearValidate()
      })
    },
    updateData() {
      this.$refs['dataForm'].validate((valid) => {
        if (valid) {
          const tempData = Object.assign({}, this.temp)
          update(tempData).then(() => {
            this.dialogFormVisible = false
            this.$notify({
              title: '成功',
              message: '更新成功',
              type: 'success',
              duration: 2000
            })
            this.getList()
          })
        }
      })
    },
    handleDelete(row) {
      this.$confirm('确认删除链?', '删除链', {
        confirmButtonText: this.$t('button.confirm'),
        cancelButtonText: this.$t('button.cancel'),
        type: 'warning'
      })
        .then(async() => {
          await del(row)
          this.$notify({
            title: '成功',
            message: '删除成功',
            type: 'success',
            duration: 2000
          })
          this.getList()
        })
        .catch(err => { console.error(err) })
    },
    handleBuild(row) {
      this.$confirm('确认构建链?', '构建链', {
        confirmButtonText: this.$t('button.confirm'),
        cancelButtonText: this.$t('button.cancel'),
        type: 'warning'
      })
        .then(async() => {
          await build(row)
          this.$notify({
            title: '成功',
            message: '构建成功',
            type: 'success',
            duration: 2000
          })
          this.getList()
        })
        .catch(err => { console.error(err) })
    },
    handleRun(row) {
      this.$confirm('确认运行链?', '运行链', {
        confirmButtonText: this.$t('button.confirm'),
        cancelButtonText: this.$t('button.cancel'),
        type: 'warning'
      })
        .then(async() => {
          await run(row)
          this.$notify({
            title: '成功',
            message: '运行成功',
            type: 'success',
            duration: 2000
          })
          this.getList()
        })
        .catch(err => { console.error(err) })
    },
    handleStop(row) {
      this.$confirm('确认停止链?', '停止链', {
        confirmButtonText: this.$t('button.confirm'),
        cancelButtonText: this.$t('button.cancel'),
        type: 'warning'
      })
        .then(async() => {
          await stop(row)
          this.$notify({
            title: '成功',
            message: '停止成功',
            type: 'success',
            duration: 2000
          })
          this.getList()
        })
        .catch(err => { console.error(err) })
    },
    handleRelease(row) {
      this.$confirm('确认释放链资源?', '释放链资源', {
        confirmButtonText: this.$t('button.confirm'),
        cancelButtonText: this.$t('button.cancel'),
        type: 'warning'
      })
        .then(async() => {
          await release(row)
          this.$notify({
            title: '成功',
            message: '释放成功',
            type: 'success',
            duration: 2000
          })
          this.getList()
        })
        .catch(err => { console.error(err) })
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

<style>
  .el-tag + .el-tag {
    margin-left: 5px;
  }
  .button-new-tag {
    margin-left: 10px;
    height: 32px;
    line-height: 30px;
    padding-top: 0;
    padding-bottom: 0;
  }
  .input-new-tag {
    width: 90px;
    margin-left: 10px;
    vertical-align: bottom;
  }
</style>
