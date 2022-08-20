
<script lang="ts">
	let query: string
</script>

<form class="w-full lg:w-1/2 m-4">
	<label for="search-bar" class="sr-only">Search</label>
	<div class="relative">
		<button class="flex absolute inset-y-0 pl-3 items-center">
			<svg aria-hidden="true" class="w-5 h-5 text-gray-500 dark:text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path></svg>
		</button>
		<input id="search-bar" type="search" autofocus
			class="w-full p-2 pl-10 text-lg border border-gray-300 rounded-lg focus:ring-blue"
			placeholder="Search..."
			bind:value={query}
		/>
	</div>

	{#if query}
		<p>
			Found results in time s.
		</p>
		<div class="space-x-4">
			Result 1
		</div>
	{/if}

</form>
