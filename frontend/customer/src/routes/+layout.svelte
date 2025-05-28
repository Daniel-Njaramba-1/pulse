<script lang="ts">
	import { page } from '$app/state';
  	import Header from '$lib/resources/header.svelte';
	import '../app.css';
	let { children } = $props();
</script>

{#if page.url.pathname.startsWith('/register') || page.url.pathname.startsWith('/login')}
	<main class="">
		{@render children()}
	</main>
{:else}
	<Header />
	<main class="">
		{@render children()}
	</main>
{/if}
