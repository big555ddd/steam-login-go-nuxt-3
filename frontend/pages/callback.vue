<template>
    <div class="flex items-center justify-center min-h-screen bg-gray-100">
      <div class="bg-white p-8 rounded shadow-md text-center">
        <h1 class="text-2xl font-bold mb-4 ">Login Successful</h1>
  
        <p><strong>Token:</strong></p>
        <textarea 
          readonly 
          class="w-full p-2 mt-2 bg-gray-200 rounded" 
          rows="4"
        >{{ token }}</textarea>
  
        <div class="mt-8">
          <h2 class="text-xl font-semibold">User Profile</h2>
          <p><strong>Steam ID:</strong> {{ profile.steam_id }}</p>
          <p><strong>Username:</strong> {{ profile.persona_name }}</p>
  
          <div class="flex space-x-5 justify-center">
            <img 
              :src="profile.avatar_small" 
              alt="Avatar Small" 
              class="rounded-full mt-4"
            />
            <img 
              :src="profile.avatar_medium" 
              alt="Avatar Medium" 
              class="rounded-full mt-4"
            />
            <img 
              :src="profile.avatar_full" 
              alt="Avatar Full" 
              class="rounded-full mt-4"
            />
          </div>
        </div>
  
        <button 
          class="mt-6 px-4 py-2 bg-red-500 text-white rounded"
          @click="logout"
        >
          Logout
        </button>
      </div>
    </div>
  </template>
  

<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'
import { ref, onMounted } from 'vue'
import Cookies from 'js-cookie'

// ตัวแปร token และ profile
const token = ref<string | null>(null)
interface UserProfile {
  steamid: string
  personaname: string
  avatarfull: string
  profileurl: string
}
const profile = ref<UserProfile>({
  steamid: '',
  personaname: '',
  avatarfull: '',
  profileurl: '',
})

const route = useRoute()
const router = useRouter()

// ฟังก์ชัน Logout ลบ JWT จาก Cookie และกลับหน้าแรก
const logout = () => {
  Cookies.remove('jwt_token')
  router.push('/')
}

onMounted(async () => {
  const steamID = new URLSearchParams(window.location.search).get('openid.identity')?.split('/').pop()

  if (steamID) {
    try {
      const response = await fetch(`http://localhost:8080/auth/steam/${steamID}`)
      if (!response.ok) throw new Error('Failed to fetch Steam profile')

      const data = await response.json()
      token.value = data.token
      profile.value = data.profile

      // เก็บ JWT Token ลงใน Cookie มีอายุ 1 วัน
      Cookies.set('jwt_token', data.token, { expires: 1 })
    } catch (error) {
      console.error('Login failed:', error)
      router.push('/') // กลับหน้าแรกถ้าล้มเหลว
    }
  } else {
    console.error('Failed to get Steam ID')
    router.push('/') // กลับหน้าแรกถ้าไม่มี Steam ID
  }
})
</script>
