{{define "list_metrics"}}
		<html>
			<body>
				<h1>Metric List</h1>
				<h2>Gauges:</h2>
				<ul>
					{{range $name, $value := .}}
					{{if eq $value.MType "gauge"}}
						<li><strong>{{ $value.ID }}:</strong> {{ $value.Value }}</li>
							{{end}}
					{{end}}
				</ul>
				<h2>Counters:</h2>
				<ul>
					{{range $name, $value := .}}
					{{if eq $value.MType "counter"}}
						<li><strong>{{ $value.ID }}:</strong> {{ $value.Delta }}</li>
						{{end}}
					{{end}}
				</ul>
			</body>
		</html>
{{end}}
