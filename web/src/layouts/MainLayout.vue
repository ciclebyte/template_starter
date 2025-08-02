<template>
  <NavBar />
  <div class="main-content">
    <Container>
      <transition
        mode="out-in"
        @before-enter="beforeEnter"
        @enter="enter"
        @leave="leave"
      >
        <router-view />
      </transition>
    </Container>
  </div>
  <FooterBar />
</template>

<script setup>
import NavBar from '@/components/NavBar.vue'
import FooterBar from '@/components/FooterBar.vue'
import Container from '@/components/Container.vue'

function beforeEnter(el) {
  el.style.opacity = 0
  el.style.transform = 'translateY(30px)'
}

function enter(el, done) {
  el.offsetHeight // trigger reflow
  el.style.transition = 'opacity 0.4s ease, transform 0.4s ease'
  el.style.opacity = 1
  el.style.transform = 'translateY(0)'
  setTimeout(done, 400)
}

function leave(el, done) {
  el.style.transition = 'opacity 0.3s ease, transform 0.3s ease'
  el.style.opacity = 0
  el.style.transform = 'translateY(-30px)'
  setTimeout(done, 300)
}
</script>

<style>
/* 全局样式 - 确保始终显示滚动条 */
html {
  overflow-y: scroll;
}
</style>

<style scoped>
.main-content {
  min-height: calc(100vh - 64px);
  background: #f5f5f5;
  width: 100%;
  padding-top: 64px;
  padding-bottom: 20px;
}
</style> 