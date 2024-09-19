<template>
  <NavBar/>
  <RouterView/>
  <Notifications position="top center" style="margin-top: 64px;"/>
</template>


<script setup>
import { RouterView } from 'vue-router'
import NavBar from "@/components/NavBar.vue";
import {Notifications} from "@kyvg/vue3-notification";
</script>

<style lang="scss">
@import "@/assets/bulma/bulma.scss";
@import "@/assets/fontawesome-free-6.6.0-web/css/all.css";
</style>