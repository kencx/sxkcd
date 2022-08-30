<script lang="ts" context="module">
	export type Comic = {
		id: number;
		title: string;
		date: number;
		num: number;
		img_url: string;
		alt?: string;
	};
</script>

<script lang="ts">
	export let result: Comic;

	// handle no alt text
	if (result.alt == undefined) {
		result.alt = "";
	}

	function parseDate(epoch: number): string {
		var date = new Date(epoch * 1000);
		var year = date.getFullYear();
		var month = String(date.getMonth() + 1).padStart(2, "0");
		var day = String(date.getDate()).padStart(2, "0");
		return `${year}-${month}-${day}`;
	}
</script>

<div class="result-container">
	<hgroup class="titles">
		<h4>{result.title}</h4>
		<small>
			<span>{parseDate(result.date)}</span> |
			<span><a href="https://xkcd.com/{result.num}" target="_blank">#{result.num}</a></span> |
			<span><a href="https://explainxkcd.com/{result.num}" target="_blank">explain</a></span>
		</small>
	</hgroup>
	<div class="img">
		<img src={result.img_url} alt="Comic {result.num}" loading="lazy"/>
	</div>
	<div class="alt">
		<p>{result.alt}</p>
	</div>
	<hr style="border: 0; height: 1px;">
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
