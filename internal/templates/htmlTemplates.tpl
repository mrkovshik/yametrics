{{define "list_metrics"}}
		<html>
			<body>
				<h1>Metric List</h1>
				<h2>Gauges:</h2>
				<ul>
					{{range $name, $value := .Gauges}}
						<li><strong>{{ $name }}:</strong> {{ $value }}</li>
					{{end}}
				</ul>
				<h2>Counters:</h2>
				<ul>
					{{range $name, $value := .Counters}}
						<li><strong>{{ $name }}:</strong> {{ $value }}</li>
					{{end}}
				</ul>
			</body>
		</html>
{{end}}
