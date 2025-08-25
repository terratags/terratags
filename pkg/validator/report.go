package validator

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/terratags/terratags/pkg/config"
)

// UnifiedReportData combines both direct and module resource data
type UnifiedReportData struct {
	GeneratedTime        string
	Stats                TagComplianceStats
	NonCompliantCount    int
	TotalExemptResources int
	CompliancePercentage float64
	RequiredTags         []string
	Violations           []TagViolation
	HasExcludedResources bool
	// Module-specific fields
	ModuleResources      []ModuleResourceValidation
	ModuleViolations     []TagViolation
	HasModuleResources   bool
}

// GenerateUnifiedHTMLReport generates a single report that handles both direct and module resources
func GenerateUnifiedHTMLReport(violations []TagViolation, stats TagComplianceStats, cfg *config.Config, moduleResources ...[]ModuleResourceValidation) string {
	// Calculate compliance percentage
	compliancePercentage := 0.0
	if stats.TotalResources > 0 {
		compliancePercentage = float64(stats.CompliantResources) / float64(stats.TotalResources) * 100
	}

	// Handle optional module resources
	var moduleRes []ModuleResourceValidation
	if len(moduleResources) > 0 {
		moduleRes = moduleResources[0]
	}

	// Separate direct and module violations
	var directViolations, moduleViolations []TagViolation
	for _, v := range violations {
		// Check if this is a module resource by looking for module path indicators
		if strings.Contains(v.ResourcePath, "module.") {
			moduleViolations = append(moduleViolations, v)
		} else {
			directViolations = append(directViolations, v)
		}
	}

	data := UnifiedReportData{
		GeneratedTime:        time.Now().Format("2006-01-02 15:04:05"),
		Stats:                stats,
		NonCompliantCount:    stats.TotalResources - stats.CompliantResources - stats.FullyExemptResources - stats.PartiallyExemptResources,
		TotalExemptResources: stats.FullyExemptResources + stats.PartiallyExemptResources,
		CompliancePercentage: compliancePercentage,
		RequiredTags:         cfg.Required,
		Violations:           directViolations,
		HasExcludedResources: len(stats.ExcludedAWSCCResources) > 0,
		ModuleResources:      moduleRes,
		HasModuleResources:   len(moduleRes) > 0 || len(moduleViolations) > 0,
		ModuleViolations:     moduleViolations,
	}

	tmpl, err := template.New("report").Funcs(template.FuncMap{
		"join": strings.Join,
		"add": func(a, b int) int {
			return a + b
		},
	}).Parse(getUnifiedTemplate())

	if err != nil {
		return fmt.Sprintf("Error parsing template: %v", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Sprintf("Error executing template: %v", err)
	}

	return buf.String()
}

func getUnifiedTemplate() string {
	return `<!DOCTYPE html>
<html>
<head>
    <title>Terraform Tag Compliance Report</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/js/bootstrap.bundle.min.js"></script>
    <style>
        .header-logo { max-height: 60px; margin-right: 15px; }
        .github-link { margin-left: auto; text-decoration: none; }
        .progress { height: 30px; }
        .exempt-tag { color: #fd7e14; }
        .logo-svg { height: 60px; width: 60px; margin-right: 15px; }
        .module-section { background-color: #e3f2fd; border-left: 4px solid #2196f3; }
        .module-path { font-family: monospace; color: #666; font-size: 0.9em; }
    </style>
</head>
<body class="bg-light">
    <div class="container py-4">
        <!-- Header -->
        <div class="d-flex align-items-center mb-4">
            <svg class="logo-svg" viewBox="0 0 200 200" xmlns="http://www.w3.org/2000/svg">
              <path d="M120,90 C130,75 150,75 160,85 C170,75 190,80 190,100 C190,115 175,120 160,120 
                       C160,125 155,130 145,130 C135,130 90,130 80,130 C65,130 50,120 50,105 
                       C50,90 65,80 80,85 C85,70 105,70 120,90 Z" 
                    fill="#FFE0B2" stroke="#FF8C00" stroke-width="3"/>
              <path d="M100,100 L115,100 L115,115 L107.5,122 L100,115 Z" 
                    fill="#FF8C00" stroke="#FF8C00" stroke-width="1"/>
              <circle cx="107.5" cy="105" r="2" fill="#FFE0B2"/>
              <circle cx="70" cy="110" r="35" fill="none" stroke="#FF8C00" stroke-width="6"/>
              <line x1="95" y1="135" x2="120" y2="160" stroke="#FF8C00" stroke-width="10" stroke-linecap="round"/>
            </svg>
            <h1 class="mb-0">Terraform Tag Compliance Report</h1>
            <a href="https://github.com/terratags/terratags" class="github-link" target="_blank">
                <img src="https://github.githubassets.com/images/modules/logos_page/GitHub-Mark.png" alt="GitHub" width="32">
            </a>
        </div>
        
        <p class="text-muted">Generated on: {{.GeneratedTime}}</p>
        
        <!-- Summary Card -->
        <div class="card mb-4">
            <div class="card-header bg-primary text-white">
                <h2 class="card-title h5 mb-0">Summary</h2>
            </div>
            <div class="card-body">
                <div class="row">
                    <div class="col-md-3">
                        <div class="card text-center mb-3">
                            <div class="card-body">
                                <h3 class="h2">{{.Stats.TotalResources}}</h3>
                                <p class="mb-0">Total Resources</p>
                            </div>
                        </div>
                    </div>
                    <div class="col-md-2">
                        <div class="card text-center mb-3 bg-success text-white">
                            <div class="card-body">
                                <h3 class="h2">{{.Stats.CompliantResources}}</h3>
                                <p class="mb-0">Compliant</p>
                            </div>
                        </div>
                    </div>
                    <div class="col-md-2">
                        <div class="card text-center mb-3 bg-danger text-white">
                            <div class="card-body">
                                <h3 class="h2">{{.NonCompliantCount}}</h3>
                                <p class="mb-0">Non-compliant</p>
                            </div>
                        </div>
                    </div>
                    <div class="col-md-2">
                        <div class="card text-center mb-3 bg-warning">
                            <div class="card-body">
                                <h3 class="h2">{{.TotalExemptResources}}</h3>
                                <p class="mb-0">Exempt</p>
                            </div>
                        </div>
                    </div>
                    <div class="col-md-3">
                        <div class="card text-center mb-3 bg-info text-white">
                            <div class="card-body">
                                <h3 class="h2">{{.Stats.ExcludedResourcesCount}}</h3>
                                <p class="mb-0">Excluded</p>
                            </div>
                        </div>
                    </div>
                </div>
                
                <div class="progress mt-3">
                    <div class="progress-bar bg-success" role="progressbar" 
                         style="width: {{printf "%.1f" .CompliancePercentage}}%;" 
                         aria-valuenow="{{printf "%.1f" .CompliancePercentage}}" 
                         aria-valuemin="0" aria-valuemax="100">
                        {{printf "%.1f" .CompliancePercentage}}% Compliant
                    </div>
                </div>
            </div>
        </div>
        
        <!-- Direct Resources -->
        {{if .Violations}}
        <div class="card mb-4">
            <div class="card-header bg-secondary text-white">
                <h2 class="card-title h5 mb-0">Direct Resources</h2>
            </div>
            <div class="card-body">
                <div class="accordion" id="directResourceAccordion">
                    {{range $index, $v := .Violations}}
                    <div class="accordion-item">
                        <h2 class="accordion-header">
                            <button class="accordion-button collapsed" type="button" 
                                    data-bs-toggle="collapse" data-bs-target="#direct{{$index}}">
                                {{$v.ResourceType}} "{{$v.ResourceName}}"
                                {{if $v.IsExempt}}<span class="badge bg-warning ms-2">EXEMPT</span>{{end}}
                            </button>
                        </h2>
                        <div id="direct{{$index}}" class="accordion-collapse collapse">
                            <div class="accordion-body">
                                <p><strong>Path:</strong> {{$v.ResourcePath}}</p>
                                {{if $v.MissingTags}}<p><strong>Missing:</strong> {{join $v.MissingTags ", "}}</p>{{end}}
                                {{if $v.PatternViolations}}
                                <p><strong>Pattern Violations:</strong></p>
                                <ul>{{range $v.PatternViolations}}<li>{{.TagName}}: {{.ErrorMessage}}</li>{{end}}</ul>
                                {{end}}
                            </div>
                        </div>
                    </div>
                    {{end}}
                </div>
            </div>
        </div>
        {{end}}
        
        <!-- Module Resources -->
        {{if .HasModuleResources}}
        <div class="card mb-4 module-section">
            <div class="card-header bg-info text-white">
                <h2 class="card-title h5 mb-0">Module Resources</h2>
            </div>
            <div class="card-body">
                <div class="accordion" id="moduleResourceAccordion">
                    {{range $index, $m := .ModuleResources}}
                    <div class="accordion-item">
                        <h2 class="accordion-header">
                            <button class="accordion-button collapsed" type="button" 
                                    data-bs-toggle="collapse" data-bs-target="#module{{$index}}">
                                {{$m.Type}} "{{$m.Name}}"
                                <span class="module-path ms-2">({{$m.ModulePath}})</span>
                            </button>
                        </h2>
                        <div id="module{{$index}}" class="accordion-collapse collapse">
                            <div class="accordion-body">
                                <p><strong>Module:</strong> {{$m.ModulePath}} ({{$m.ModuleSource}})</p>
                                {{if $m.MissingTags}}<p><strong>Missing:</strong> {{join $m.MissingTags ", "}}</p>{{end}}
                                {{if $m.PatternViolations}}
                                <p><strong>Pattern Violations:</strong></p>
                                <ul>{{range $m.PatternViolations}}<li>{{.TagName}}: {{.Error}}</li>{{end}}</ul>
                                {{end}}
                            </div>
                        </div>
                    </div>
                    {{end}}
                    {{range $index, $v := .ModuleViolations}}
                    <div class="accordion-item">
                        <h2 class="accordion-header">
                            <button class="accordion-button collapsed" type="button" 
                                    data-bs-toggle="collapse" data-bs-target="#moduleViol{{$index}}">
                                {{$v.ResourceType}} "{{$v.ResourceName}}"
                                <span class="module-path ms-2">({{$v.ResourcePath}})</span>
                            </button>
                        </h2>
                        <div id="moduleViol{{$index}}" class="accordion-collapse collapse">
                            <div class="accordion-body">
                                <p><strong>Module Path:</strong> {{$v.ResourcePath}}</p>
                                {{if $v.MissingTags}}<p><strong>Missing:</strong> {{join $v.MissingTags ", "}}</p>{{end}}
                                {{if $v.PatternViolations}}
                                <p><strong>Pattern Violations:</strong></p>
                                <ul>{{range $v.PatternViolations}}<li>{{.TagName}}: {{.ErrorMessage}}</li>{{end}}</ul>
                                {{end}}
                            </div>
                        </div>
                    </div>
                    {{end}}
                </div>
            </div>
        </div>
        {{end}}
        
        <footer class="mt-4 text-center text-muted">
            <p>Generated by <a href="https://github.com/terratags/terratags" target="_blank">Terratags</a></p>
        </footer>
    </div>
</body>
</html>`
}
