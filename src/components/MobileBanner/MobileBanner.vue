<template>
  <div v-if="isMobileOrTablet" class="mobile-banner">
    This version of the breaking news demo is not usable on mobile/tablet devices.
    Please use this demo on a computer.
  </div>
</template>

<script>
export default {
  name: "MobileBanner",
  computed: {
    isMobileOrTablet() {
      const userAgent = navigator.userAgent.toLowerCase();

      if (
        userAgent.includes("android") ||
        userAgent.includes("iphone") ||
        userAgent.includes("ipad") ||
        userAgent.includes("mobile") ||
        window.matchMedia("(max-width: 1080px)").matches
      ) {
        return true;
      } 
       
      return false
    },
  },
};
</script>

<style>
.mobile-banner {
  background-color: red;
  color: white;
  text-align: center;
  padding: 20px;
  font-size: 16px;
}
</style>

