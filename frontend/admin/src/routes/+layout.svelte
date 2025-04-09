<script lang="ts">
  	import { page } from '$app/state';
  	import Sidebar from '$lib/resources/sidebar.svelte';
	import '../app.css';
	let { children } = $props();
</script>

{#if page.url.pathname.startsWith('/register') || page.url.pathname.startsWith('/login')}
	<main class="">
		{@render children()}
	</main>
{:else}
	<Sidebar />
	<main class="ml-[60px]">
		{@render children()}
	</main>
{/if}
