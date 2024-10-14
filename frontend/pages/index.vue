<template>
    <div class="flex items-center justify-center min-h-screen bg-gray-100">
      <div class="bg-white p-8 rounded shadow-md text-center">
        <h1 class="text-2xl font-bold mb-4">Login with Steam</h1>
        <button
          class="bg-green-500 hover:bg-green-600 text-white font-semibold py-2 px-4 rounded"
          @click="loginWithSteam"
        >
          Login with Steam
        </button>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  function loginWithSteam() {
    const steamLoginURL = `https://steamcommunity.com/openid/login?openid.ns=http://specs.openid.net/auth/2.0&openid.mode=checkid_setup&openid.return_to=http://localhost:3000/callback&openid.realm=http://localhost:3000&openid.identity=http://specs.openid.net/auth/2.0/identifier_select&openid.claimed_id=http://specs.openid.net/auth/2.0/identifier_select`
    window.location.href = steamLoginURL // Redirect ไป Steam เพื่อ Login
  }
  </script>
  