<template>
  <div id="userLayout" class="w-full h-full relative">
    <div
      class="rounded-lg flex items-center justify-evenly w-full h-full md:w-screen md:h-screen md:bg-[#194bfb] bg-white"
    >
      <div class="md:w-3/5 w-10/12 h-full flex items-center justify-evenly">
        <div
          class="oblique h-[130%] w-3/5 bg-white dark:bg-slate-900 transform -rotate-12 absolute -ml-52"
        />
        <!-- 分割斜块 -->
        <div
          class="z-[999] pt-12 pb-10 md:w-96 w-full rounded-lg flex flex-col justify-between box-border"
        >
          <div>
            <div class="flex items-center justify-center">
              <img class="w-24" :src="$GIN_VUE_ADMIN.appLogo" alt />
            </div>
            <div class="mb-9">
              <p class="text-center text-4xl font-bold">
                用户注册
              </p>
              <p class="text-center text-sm font-normal text-gray-500 mt-2.5">
                A medical privacy sharing and protection system based on Heyperledger Fabric
              </p>
            </div>
            <el-form
              ref="registerForm"
              :model="registerFormData"
              :rules="rules"
              :validate-on-rule-change="false"
              @keyup.enter="submitForm"
            >
              <el-form-item prop="username" class="mb-6">
                <el-input
                  v-model="registerFormData.username"
                  size="large"
                  placeholder="请输入用户名"
                  suffix-icon="user"
                />
              </el-form-item>
              <el-form-item prop="password" class="mb-6">
                <el-input
                  v-model="registerFormData.password"
                  show-password
                  size="large"
                  type="password"
                  placeholder="请输入密码"
                />
              </el-form-item>
              <el-form-item prop="confirmPassword" class="mb-6">
                <el-input
                  v-model="registerFormData.confirmPassword"
                  show-password
                  size="large"
                  type="password"
                  placeholder="请确认密码"
                />
              </el-form-item>
              <el-form-item prop="realName" class="mb-6">
                <el-input
                  v-model="registerFormData.realName"
                  size="large"
                  placeholder="请输入真实姓名"
                />
              </el-form-item>
              <el-form-item prop="idCard" class="mb-6">
                <el-input
                  v-model="registerFormData.idCard"
                  size="large"
                  placeholder="请输入身份证号"
                />
              </el-form-item>
              <el-form-item prop="phone" class="mb-6">
                <el-input
                  v-model="registerFormData.phone"
                  size="large"
                  placeholder="请输入手机号"
                />
              </el-form-item>
              <el-form-item prop="smsCode" class="mb-6">
                <div class="flex w-full justify-between">
                  <el-input
                    v-model="registerFormData.smsCode"
                    placeholder="请输入短信验证码"
                    size="large"
                    class="flex-1 mr-5"
                  />
                  <el-button
                    type="primary"
                    size="large"
                    :disabled="smsCodeBtnDisabled"
                    @click="handleSendSmsCode"
                  >{{ smsCodeBtnText }}</el-button>
                </div>
              </el-form-item>
              <!-- <el-form-item
                v-if="registerFormData.openCaptcha"
                prop="captcha"
                class="mb-6"
              >
                <div class="flex w-full justify-between">
                  <el-input
                    v-model="registerFormData.captcha"
                    placeholder="请输入图形验证码"
                    size="large"
                    class="flex-1 mr-5"
                  />
                  <div class="w-1/3 h-11 bg-[#c3d4f2] rounded">
                    <img
                      v-if="picPath"
                      class="w-full h-full"
                      :src="picPath"
                      alt="请输入验证码"
                      @click="getImageCaptcha()"
                    />
                  </div>
                </div>
              </el-form-item> -->
              <el-form-item class="mb-6">
                <el-button
                  class="shadow shadow-active h-11 w-full"
                  type="primary"
                  size="large"
                  @click.prevent="submitForm"
                  >注 册</el-button
                >
              </el-form-item>
              <el-form-item class="mb-6">
                <el-button
                  class="shadow shadow-active h-11 w-full"
                  size="large"
                  @click="goToLogin"
                  >返回登录</el-button
                >
              </el-form-item>
            </el-form>
          </div>
        </div>
      </div>
    </div>

    <!-- <BottomInfo class="left-0 right-0 absolute bottom-3 mx-auto w-full z-20">
      <div class="links items-center justify-center gap-2 hidden md:flex">
        <a href="https://www.gin-vue-admin.com/" target="_blank">
          <img src="@/assets/docs.png" class="w-8 h-8" alt="文档" />
        </a>
        <a href="https://support.qq.com/product/371961" target="_blank">
          <img src="@/assets/kefu.png" class="w-8 h-8" alt="客服" />
        </a>
        <a
          href="https://github.com/flipped-aurora/gin-vue-admin"
          target="_blank"
        >
          <img src="@/assets/github.png" class="w-8 h-8" alt="github" />
        </a>
        <a href="https://space.bilibili.com/322210472" target="_blank">
          <img src="@/assets/video.png" class="w-8 h-8" alt="视频站" />
        </a>
      </div>
    </BottomInfo> -->
  </div>
</template>

<script setup>
import { captcha, register, sendSmsCode } from '@/api/user'
import BottomInfo from '@/components/bottomInfo/bottomInfo.vue'
import { reactive, ref,onBeforeUnmount } from 'vue'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'

defineOptions({
  name: 'Register'
})

const router = useRouter()

// 验证函数
const checkUsername = (rule, value, callback) => {
  if (value.length < 5) {
    return callback(new Error('用户名不能少于5个字符'))
  } else {
    callback()
  }
}

const checkPassword = (rule, value, callback) => {
  if (value.length < 6) {
    return callback(new Error('密码不能少于6个字符'))
  } else {
    callback()
  }
}

const checkConfirmPassword = (rule, value, callback) => {
  if (value !== registerFormData.password) {
    return callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const checkRealName = (rule, value, callback) => {
  if (!value) {
    return callback(new Error('请输入真实姓名'))
  } else {
    callback()
  }
}

const checkIdCard = (rule, value, callback) => {
  const idCardReg = /(^\d{15}$)|(^\d{18}$)|(^\d{17}(\d|X|x)$)/
  if (!idCardReg.test(value)) {
    return callback(new Error('请输入正确的身份证号'))
  } else {
    callback()
  }
}

const checkPhone = (rule, value, callback) => {
  const phoneReg = /^1[3-9]\d{9}$/
  if (!phoneReg.test(value)) {
    return callback(new Error('请输入正确的手机号'))
  } else {
    callback()
  }
}

const checkSmsCode = (rule, value, callback) => {
  if (!value || value.length !== 6) {
    return callback(new Error('请输入6位短信验证码'))
  } else {
    callback()
  }
}

// // 获取图形验证码
// const getImageCaptcha = async () => {
//   const ele = await captcha()
//   rules.captcha.push({
//     max: ele.data.captchaLength,
//     min: ele.data.captchaLength,
//     message: `请输入${ele.data.captchaLength}位验证码`,
//     trigger: 'blur'
//   })
//   picPath.value = ele.data.picPath
//   registerFormData.captchaId = ele.data.captchaId
//   registerFormData.openCaptcha = ele.data.openCaptcha
// }
// getImageCaptcha()

// 注册相关操作
const registerForm = ref(null)
const picPath = ref('')
const registerFormData = reactive({
  username: '',
  password: '',
  confirmPassword: '',
  realName: '',
  idCard: '',
  phone: '',
  smsCode: '',
  captcha: '',
  captchaId: '',
  openCaptcha: false
})

const rules = reactive({
  username: [{ validator: checkUsername, trigger: 'blur' }],
  password: [{ validator: checkPassword, trigger: 'blur' }],
  confirmPassword: [{ validator: checkConfirmPassword, trigger: 'blur' }],
  realName: [{ validator: checkRealName, trigger: 'blur' }],
  idCard: [{ validator: checkIdCard, trigger: 'blur' }],
  phone: [{ validator: checkPhone, trigger: 'blur' }],
  smsCode: [{ validator: checkSmsCode, trigger: 'blur' }],
  captcha: [
    {
      message: '验证码格式不正确',
      trigger: 'blur'
    }
  ]
})

// 短信验证码相关逻辑
const smsCodeBtnText = ref('发送验证码')
const smsCodeBtnDisabled = ref(false)
let countdownTimer = null

// 更新发送短信验证码的函数
const handleSendSmsCode = async () => {
  try {
    const isPhoneValid = await validateField('phone')
    if (!isPhoneValid) return
    
    // 调用实际API
    await sendSmsCode(registerFormData.phone)
    
    ElMessage({
      type: 'success',
      message: '验证码已发送',
      showClose: true
    })
    
    startCountdown()
  } catch (error) {
    ElMessage({
      type: 'error',
      message: error.message || '发送验证码失败',
      showClose: true
    })
  }
}

// 验证单个字段
const validateField = async (field) => {
  try {
    await registerForm.value.validateField(field)
    return true
  } catch (error) {
    return false
  }
}

// 开始倒计时
const startCountdown = () => {
  let countdown = 60
  smsCodeBtnDisabled.value = true
  smsCodeBtnText.value = `重新发送(${countdown}s)`
  
  countdownTimer = setInterval(() => {
    countdown--
    smsCodeBtnText.value = `重新发送(${countdown}s)`
    
    if (countdown <= 0) {
      clearInterval(countdownTimer)
      smsCodeBtnText.value = '发送验证码'
      smsCodeBtnDisabled.value = false
    }
  }, 1000)
}

// 更新提交注册表单的函数
const submitForm = () => {
    console.log('提交按钮被点击')  // 添加这行
    if (!registerForm.value) {
    console.error('表单引用无效')
    ElMessage({
      type: 'error',
      message: '表单引用无效',
      showClose: true
    })
    return
  }



  registerForm.value.validate(async (v) => {
    console.log('表单验证结果:', v)  // 添加这行
    if (!v) {
      ElMessage({
        type: 'error',
        message: '请正确填写注册信息',
        showClose: true
      })
      await getImageCaptcha()
      return false
    }

    try {
      // 调用实际的注册API
      const res = await register(registerFormData)
      
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: '注册成功，请登录',
          showClose: true
        })
        
        setTimeout(() => {
          router.push({ name: 'Login' })
        }, 1500)
        
        return true
      } else {
        throw new Error(res.msg || '注册失败')
      }
    } catch (error) {
      ElMessage({
        type: 'error',
        message: error.message || '注册失败，请重试',
        showClose: true
      })
      await getImageCaptcha()
      return false
    }
  })
}

// 返回登录页
const goToLogin = () => {
  router.push({ name: 'Login' })
}

// 组件销毁时清除定时器
onBeforeUnmount(() => {
  if (countdownTimer) {
    clearInterval(countdownTimer)
  }
})
</script>