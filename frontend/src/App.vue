<script>

import Home from './components/home.vue';

export default {
  name: 'App',
  components: {
    Home,
  },
};

</script>

<template>
  <div id="main"></div>
  <router-view></router-view>
</template>

<style>
#app{
  height: 100vh;
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: center;
}
</style>
