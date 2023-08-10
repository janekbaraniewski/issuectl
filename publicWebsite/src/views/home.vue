<template>
  <div class="home-container">
    <h1 class="home-text">issuectl</h1>
    <span class="home-text1">

    </span>
  </div>
</template>

<script>
export default {
  name: 'Home',
  metaInfo: {
    title: 'issuectl',
    meta: [
      {
        property: 'og:title',
        content: 'issuectl',
      },
    ],
  },
}
</script>

<style scoped>
.home-container {
  width: 100%;
  display: flex;
  overflow: auto;
  min-height: 100vh;
  align-items: center;
  flex-direction: column;
  justify-content: center;
}
.home-text {
  font-size: 60px;
  font-style: normal;
  font-weight: 100;
}
.home-text1 {
  font-style: normal;
  font-family: Roboto;
  font-weight: 100;
}
</style>
