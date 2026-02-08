package renderer

const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>{{.Title}}</title>
  <style>
    @page {
      size: A4;
      margin: 15mm;
    }
    * {
      box-sizing: border-box;
    }
    body {
      font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
      font-size: 11pt;
      line-height: 1.4;
      margin: 0;
      padding: 15mm;
      color: #333;
    }
    table {
      width: 100%;
      border-collapse: collapse;
      table-layout: fixed;
    }
    thead td {
      background: #f0f0f0;
      font-weight: bold;
      text-align: center;
      padding: 8px 12px;
      border-bottom: 2px solid #ccc;
      font-size: 12pt;
    }
    tbody tr {
      page-break-inside: avoid;
    }
    td {
      width: 50%;
      padding: 6px 12px;
      vertical-align: top;
      border-bottom: 1px solid #eee;
    }
    td:first-child {
      border-right: 1px solid #ddd;
    }
    td h1, td h2, td h3, td h4, td h5, td h6 {
      margin-top: 0.3em;
      margin-bottom: 0.2em;
    }
    td p {
      margin-top: 0.2em;
      margin-bottom: 0.2em;
    }
    td ul, td ol {
      margin-top: 0.2em;
      margin-bottom: 0.2em;
      padding-left: 1.5em;
    }
    code {
      background: #f8f8f8;
      padding: 1px 4px;
      border-radius: 3px;
      font-size: 10pt;
    }
    pre {
      background: #f8f8f8;
      padding: 8px;
      border-radius: 4px;
      overflow-x: auto;
      font-size: 9pt;
    }
    pre code {
      background: none;
      padding: 0;
    }
    blockquote {
      border-left: 3px solid #ddd;
      margin: 0.3em 0;
      padding: 0.2em 0 0.2em 1em;
      color: #666;
    }
    hr {
      border: none;
      border-top: 1px solid #ddd;
      margin: 0.5em 0;
    }
  </style>
</head>
<body>
  <table>
    <thead>
      <tr>
        <td>{{.SourceLabel}}</td>
        <td>{{.TargetLabel}}</td>
      </tr>
    </thead>
    <tbody>
      {{range .Pairs}}
      <tr>
        <td>{{.Source}}</td>
        <td>{{.Target}}</td>
      </tr>
      {{end}}
    </tbody>
  </table>
</body>
</html>`
