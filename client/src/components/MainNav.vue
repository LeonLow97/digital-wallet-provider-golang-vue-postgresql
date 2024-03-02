<template>
  <div class="nav_container">
    <div>
      <h2 class="nav_logo">Mobile Wallet</h2>
    </div>
    <div>
      <ul class="nav_list_items">
        <li><a href="">Wallets</a></li>
        <li><a href="">Make a Transfer</a></li>
        <li><a href="">Balance</a></li>
        <li><a href="">Profile</a></li>
      </ul>
    </div>
  </div>
</template>

<style scoped>
.nav_container {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-bottom: 3.2rem;
  padding-top: 2rem;
}

.nav_logo { 
  margin-left: 2.4rem;
}

.nav_list_items {
  display: flex;
  gap: 2.4rem;
  list-style: none;
  font-size: 2rem;
  margin-right: 2.4rem;
}
</style>
