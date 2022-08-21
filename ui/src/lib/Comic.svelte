<script lang="ts" context="module">
	export type Comic = {
		title: string;
		date: string;
		num: number;
		img_url: string;
		alt: string;
	};
</script>

<script lang="ts">
	export let result: Comic;
</script>

<div class="result-container">
	<hgroup class="titles">
		<h4>{result.title}</h4>
		<small>
			<span>{result.date}</span> |
			<span><a href="https://xkcd.com/{result.num}" target="_blank">#{result.num}</a></span> |
			<span><a href="https://explainxkcd.com/{result.num}" target="_blank">explain</a></span>
		</small>
	</hgroup>
	<div class="img">
		<img src={result.img_url} alt=""/>
	</div>
	<div class="alt">
		<p>{result.alt}</p>
	</div>
	<hr style="border: 0; height: 1px; background-color: #333;">
</div>

<style>
.result-container {
	display: flex;
	flex-direction: column;
	justify-content: center;
	align-items: center;
}

.titles small {
	font-size: small;
}

.img {
	margin: 1rem;
	margin-top: 0;
}

.alt {
	width: 85%;
}

.alt p {
	text-align: center;
	font-size: small;
}
</style>
