<template>
  <div class="register-container">
    <h2 class="text-center">注册</h2>
    <el-form ref="registerForm" :model="registerFormData" :rules="rules" label-width="120px">
      <el-form-item label="真实姓名" prop="realName">
        <el-input v-model="registerFormData.realName" placeholder="请输入真实姓名" />
      </el-form-item>
      <el-form-item label="身份证号" prop="idNumber">
        <el-input v-model="registerFormData.idNumber" placeholder="请输入身份证号" />
      </el-form-item>
      <el-form-item label="手机号" prop="phoneNumber">
        <el-input v-model="registerFormData.phoneNumber" placeholder="请输入手机号" />
      </el-form-item>
      <el-form-item label="验证码" prop="verificationCode">
        <div class="flex">
          <el-input v-model="registerFormData.verificationCode" placeholder="请输入验证码" />
          <el-button :disabled="isCodeSent" @click="sendVerificationCode">{{ isCodeSent ? countdown + '秒后重发' : '发送验证码' }}</el-button>
        </div>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="submitForm">注册</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup>
import { ref, reactive, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { sendCode } from '@/api/register'

const registerFormData = reactive({
  realName: '',
  idNumber: '',
  phoneNumber: '',
  verificationCode: ''
})

const rules = {
  realName: [{ required: true, message: '请输入真实姓名', trigger: 'blur' }],
  idNumber: [{ required: true, message: '请输入身份证号', trigger: 'blur' }],
  phoneNumber: [{ required: true, message: '请输入手机号', trigger: 'blur' }],
  verificationCode: [{ required: true, message: '请输入验证码', trigger: 'blur' }]
}

const isCodeSent = ref(false)
const countdown = ref(60)

const sendVerificationCode = async () => {
  if (isCodeSent.value) return

  try {
    await sendCode(registerFormData.phoneNumber)
    ElMessage({ type: 'success', message: '验证码已发送' })
    isCodeSent.value = true
    startCountdown()
  } catch (error) {
    ElMessage({ type: 'error', message: '发送验证码失败，请重试' })
  }
}

const startCountdown = () => {
  countdown.value = 60
  const timer = setInterval(() => {
    countdown.value--
    if (countdown.value <= 0) {
      clearInterval(timer)
      isCodeSent.value = false
    }
  }, 1000)
}

const submitForm = () => {
  registerForm.value.validate((valid) => {
    if (valid) {
      // Handle registration logic here
      ElMessage({ type: 'success', message: '注册成功' })
    } else {
      ElMessage({ type: 'error', message: '请正确填写注册信息' })
      return false
    }
  })
}
</script>

<style scoped>
.register-container {
  max-width: 400px;
  margin: auto;
}
</style>