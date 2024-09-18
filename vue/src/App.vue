<template>
  <NavBar/>

  <RouterView/>
</template>


<script setup>
import { RouterView } from 'vue-router'
import NavBar from "@/components/NavBar.vue";
</script>

<style lang="scss">
@import "@/assets/bulma/bulma.scss";
@import "@/assets/fontawesome-free-6.6.0-web/css/all.css";
</style>