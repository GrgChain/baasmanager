<template>
  <div class="app-container">
    <div class="filter-container">
      <el-input v-model="listQuery.chaincodeName" :placeholder="$t('chaincode.chaincodeName')" style="width: 200px;" class="filter-item" @keyup.enter.native="handleFilter" />
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
      <el-table-column :label="$t('chaincode.id')" prop="id" align="center" width="80">
        <template slot-scope="scope">
          <span>{{ scope.row.id }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('chaincode.chaincodeName')" min-width="80px">
        <template slot-scope="{row}">
          <span>{{ row.chaincodeName }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('chaincode.version')" min-width="40px">
        <template slot-scope="{row}">
          <span>{{ row.version }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('chaincode.created')" min-width="80px">
        <template slot-scope="{row}">
          <span>{{ row.created | parseTime('{y}-{m}-{d} {h}:{i}') }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('chaincode.status')" min-width="40px">
        <template slot-scope="{row}">
          <el-tag v-if="row.status==0" type="info">定义</el-tag>
          <el-tag v-if="row.status==1" type="warning">运行中</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="$t('chaincode.policy')" min-width="80px">
        <template slot-scope="{row}">
          <span>{{ row.policy }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('chaincode.args')" min-width="120px">
        <template slot-scope="{row}">
          <span>{{ row.args }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('button.actions')" align="center" width="230" class-name="small-padding fixed-width">
        <template slot-scope="{row}">
          <el-button v-if="row.status==0" type="primary" size="mini" :loading="isload" @click="handleDeploy(row)">
            {{ $t('button.deploy') }}
          </el-button>
          <el-button v-if="row.status==1" type="success" size="mini" :loading="isload" @click="handleUpdate(row)">
            {{ $t('button.upgrade') }}
          </el-button>
          <el-button v-if="row.status==1" type="warning" size="mini" @click="handleOperation(row)">
            {{ $t('button.operation') }}
          </el-button>
          <el-button v-if="row.status==0" size="mini" type="danger" @click="handleDelete(row)">
            {{ $t('button.delete') }}
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <pagination v-show="total>0" :total="total" :page.sync="listQuery.page" :limit.sync="listQuery.limit" @pagination="getList" />

    <el-dialog :title="textMap[dialogStatus]" :visible.sync="dialogFormVisible">
      <el-form ref="dataForm" :rules="rules" :model="temp" label-position="left" label-width="120px" style="width: 500px; margin-left:50px;">

        <el-form-item :label="$t('chaincode.chaincodeName')" prop="chaincodeName">
          <el-input v-model="temp.chaincodeName" :disabled="isupgrade" />
        </el-form-item>
        <el-form-item :label="$t('chaincode.code')" prop="githubPath">
          <el-upload
            class="upload-demo"
            :action="uploadUrl"
            :on-remove="handleRemove"
            :before-remove="beforeRemove"
            :on-success="handleSuccess"
            multiple
            :limit="1"
            :on-exceed="handleExceed"
            :file-list="fileList"
          >
            <el-button size="small" type="primary">点击上传</el-button>
            <div slot="tip" class="el-upload__tip">只能上传go文件</div>
          </el-upload>
        </el-form-item>
        <el-form-item :label="$t('chaincode.policy')" prop="policy">
          <el-input v-model="temp.policy" />
          <span style="color: #E74C3C;">MSP:{{ msp }}</span>
          <a class="el-upload__tip" target="_blank" href="https://hyperledger-fabric.readthedocs.io/en/latest/endorsement-policies.html">查看文档</a>
        </el-form-item>
        <el-form-item :label="$t('chaincode.initArgs')">
          <el-tag v-for="tag in argTags" :key="tag" closable :disable-transitions="false" @close="handleClose(tag)">{{ tag }}</el-tag>
          <el-input v-if="inputVisible" ref="saveTagInput" v-model="inputValue" size="small" class="input-new-tag" @keyup.enter.native="handleInputConfirm" @blur="handleInputConfirm" />
          <el-button v-else class="button-new-tag" size="small" @click="showInput">+ New Arg</el-button>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible = false">
          {{ $t('button.cancel') }}
        </el-button>
        <el-button :loading="isload" type="primary" @click="dialogStatus==='create'?createData():updateData()">
          {{ $t('button.confirm') }}
        </el-button>
      </div>
    </el-dialog>

    <el-dialog :title="textMap[dialogStatus]" :visible.sync="dialogFormVisible2">
      <el-form ref="dataForm2" :rules="rules" :model="temp" label-position="left" label-width="120px" style="width: 500px; margin-left:50px;">
        <el-form-item :label="$t('chaincode.fcn')" prop="fcn">
          <el-input v-model="temp.fcn" />
        </el-form-item>
        <el-form-item :label="$t('chaincode.fcntype')" prop="fcntype">
          <el-select v-model="temp.fcntype" placeholder="请选择">
            <el-option
              v-for="item in options"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('chaincode.args')">
          <el-tag v-for="tag in argTags" :key="tag" closable :disable-transitions="false" @close="handleClose(tag)">{{ tag }}</el-tag>
          <el-input v-if="inputVisible" ref="saveTagInput" v-model="inputValue" size="small" class="input-new-tag" @keyup.enter.native="handleInputConfirm" @blur="handleInputConfirm" />
          <el-button v-else class="button-new-tag" size="small" @click="showInput">+ New Arg</el-button>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible2 = false">
          {{ $t('button.cancel') }}
        </el-button>
        <el-button :loading="isload" type="primary" @click="operaData()">
          {{ $t('button.confirm') }}
        </el-button>
      </div>
      <div v-for="t in tips" :key="t.val">
        <el-alert :title="t.val" :type="t.type" />
      </div>
    </el-dialog>

  </div>
</template>

<script>
import { fetchList, add, upgrade, del, deploy, invoke, query, queryLedger, queryLatestBlocks } from '@/api/chaincode'
import { fetch } from '@/api/channel'
import waves from '@/directive/waves' // waves directive
import Pagination from '@/components/Pagination' // secondary package based on el-pagination
import { parseTime } from '@/utils'
import { mapGetters } from 'vuex'
let vm = {}
export default {
  name: 'ComplexTable',
  components: { Pagination },
  directives: { waves },
  msp() {
    return vm.msp
  },
  filters: {
  },
  data() {
    vm = this
    return {
      uploadUrl: process.env.VUE_APP_BASE_API + '/upload',
      channelId: 0,
      channel: {
        id: 0,
        userAccount: ''
      },
      tableKey: 0,
      list: null,
      total: 0,
      listLoading: true,
      listQuery: {
        page: 1,
        limit: 10,
        chaincodeName: undefined,
        channelId: undefined
      },
      tips: [],
      options: [{
        value: 'invoke',
        label: 'Invoke'
      }, {
        value: 'query',
        label: 'Query'
      }],
      msp: '',
      isupgrade: false,
      isload: false,
      temp: {
        chaincodeName: undefined,
        channelId: 0,
        created: 0,
        version: undefined,
        githubPath: undefined,
        userAccount: undefined,
        args: undefined,
        policy: undefined
      },
      dialogFormVisible: false,
      dialogFormVisible2: false,
      dialogStatus: '',
      textMap: {
        update: '编辑链码',
        create: '创建链码',
        opera: '操作链码'
      },
      rules: {
        chaincodeName: [{ required: true, message: 'chaincodeName is required', trigger: 'blur' }],
        githubPath: [{ required: true, message: 'file is required', trigger: 'blur' }],
        policy: [{ required: true, message: 'policy is required', trigger: 'blur' }],
        fcn: [{ required: true, message: 'fcn is required', trigger: 'blur' }],
        fcntype: [{ required: true, message: 'fcntype is required', trigger: 'blur' }]
      },
      fileList: [],
      argTags: [],
      inputVisible: false,
      inputValue: ''
    }
  },
  computed: {
    ...mapGetters([
      'account'
    ])
  },
  created() {
    const id = this.$route.params && this.$route.params.id
    this.channelId = Number(id)
    this.getList()
    this.getChannel()
  },
  methods: {
    handleClose(tag) {
      this.argTags.splice(this.argTags.indexOf(tag), 1)
      this.temp.args = this.argTags.join(',')
    },
    showInput() {
      this.inputVisible = true
      this.$nextTick(_ => {
        this.$refs.saveTagInput.$refs.input.focus()
      })
    },
    handleInputConfirm() {
      const inputValue = this.inputValue.toLowerCase()
      if (inputValue) {
        this.argTags.push(inputValue)
        this.temp.args = this.argTags.join(',')
      }
      this.inputVisible = false
      this.inputValue = ''
    },
    handleRemove(file, fileList) {
      console.log(file, fileList)
      this.temp.githubPath = ''
    },
    handleSuccess(response, file, fileList) {
      console.log(response, file, fileList)
      this.temp.githubPath = response
    },
    handleExceed(files, fileList) {
      this.$message.warning('限制选择 1 个文件')
    },
    beforeRemove(file, fileList) {
      return this.$confirm('确定移除' + file.name + '？')
    },
    getChannel() {
      this.channel.id = this.channelId
      this.channel.userAccount = this.account
      fetch(this.channel).then(response => {
        if (response.code === 0) {
          const orgs = response.data.orgs.split(',')
          const msp = []
          orgs.forEach(v => {
            var o = v.toLowerCase().replace(/( |^)[a-z]/g, (L) => L.toUpperCase())
            msp.push(o + 'MSP')
          })
          this.msp = msp.join(',')
        }
      })
    },
    getList() {
      this.listLoading = true
      this.listQuery.channelId = this.channelId
      fetchList(this.listQuery).then(response => {
        this.list = response.data
        this.total = response.total
        // Just to simulate the time of the request
        setTimeout(() => {
          this.listLoading = false
        }, 1.5 * 1000)
      })

      queryLedger(this.channelId).then(response => {
        console.log(response.data)
      })
      queryLatestBlocks(this.channelId).then(response => {
        console.log(response.data)
      })
    },
    handleFilter() {
      this.listQuery.page = 1
      this.getList()
    },
    resetTemp() {
      this.fileList = []
      this.temp = {
        chaincodeName: undefined,
        channelId: 0,
        created: 0,
        version: undefined,
        githubPath: undefined,
        userAccount: undefined,
        args: undefined,
        policy: undefined
      }
    },
    handleCreate() {
      this.resetTemp()
      this.dialogStatus = 'create'
      this.dialogFormVisible = true
      this.isupgrade = false
      this.$nextTick(() => {
        this.$refs['dataForm'].clearValidate()
      })
    },
    createData() {
      this.$refs['dataForm'].validate((valid) => {
        if (valid) {
          this.temp.userAccount = this.account
          this.temp.channelId = this.channelId
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
    handleDeploy(row) {
      this.$confirm('确认部署链码?', '部署链码', {
        confirmButtonText: this.$t('button.confirm'),
        cancelButtonText: this.$t('button.cancel'),
        type: 'warning'
      })
        .then(async() => {
          this.isload = true
          await deploy(row)
          this.$notify({
            title: '成功',
            message: '部署成功',
            type: 'success',
            duration: 2000
          })
          this.isload = false
          this.getList()
        })
        .catch(err => { console.error(err) })
    },
    handleUpdate(row) {
      this.temp = Object.assign({}, row) // copy obj
      this.dialogStatus = 'update'
      this.dialogFormVisible = true
      this.isupgrade = true
      this.fileList = []
      this.argTags = this.temp.args.split(',')
      this.$nextTick(() => {
        this.$refs['dataForm'].clearValidate()
      })
    },
    updateData() {
      this.$refs['dataForm'].validate((valid) => {
        if (valid) {
          this.isload = true
          const tempData = Object.assign({}, this.temp)
          upgrade(tempData).then(() => {
            this.dialogFormVisible = false
            this.isload = false
            this.$notify({
              title: '成功',
              message: '升级成功',
              type: 'success',
              duration: 2000
            })
            this.getList()
          })
        }
      })
    },
    handleOperation(row) {
      this.temp = Object.assign({}, row) // copy obj
      this.dialogFormVisible2 = true
      this.dialogStatus = 'opera'
      this.argTags = []
      this.$nextTick(() => {
        this.$refs['dataForm2'].clearValidate()
      })
    },
    operaData() {
      this.$refs['dataForm2'].validate((valid) => {
        if (valid) {
          this.isload = true
          this.temp.args = this.temp.fcn + ',' + this.temp.args
          const tempData = Object.assign({}, this.temp)
          if (this.temp.fcntype === 'invoke') {
            invoke(tempData).then(response => {
              this.dialogFormVisible = false
              this.isload = false
              var result = {}
              result.val = response.data
              result.type = 'success'
              this.tips.push(result)
            })
          } else {
            query(tempData).then(response => {
              this.dialogFormVisible = false
              this.isload = false
              var result = {}
              result.val = response.data
              result.type = 'warning'
              this.tips.push(result)
            })
          }
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
