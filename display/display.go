package display

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"ScopeGuardian/domains/models"
	environment_variable "ScopeGuardian/environnement_variable"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

const (
	caption           = "For further results' details, please refer to the \"results\" folder located in your shared volume: %s"
	rowEngine         = "Engine"
	rowSeverity       = "Severity"
	rowName           = "Name"
	rowCwe            = "CWE"
	rowDescription    = "Description"
	rowSinkFile       = "Sink File"
	rowSinkLine       = "Sink Line"
	rowRecommendation = "Recommendation"
	rowStatus         = "Status"

	// CSV column headers use compact names (no spaces) to be machine-readable.
	csvColSinkFile = "SinkFile"
	csvColSinkLine = "SinkLine"

	errUnsupportedFormat = "unsupported format: %s"
)

// DisplayBanner prints the ASCII art banner for ScopeGuardian to w.
func DisplayBanner(w io.Writer) {
	fmt.Fprintln(w, `                                                                                   
                                                                                                @@  
                                                                                            =@@@@@  
                                                                                         @@@@@@@@   
                                                                                     =@@@@@@@@@@#   
                                                                                  @@@@@@@@@@@@@@    
                                                                              =@@@@@@@@@@@@@@@@.    
                                                                           @@@@@@@@@@@@@@@@@@@@     
                                                                       =@@@@@@@@@@@@@@@@@@@@@@      
                                                                    @@@@@@@@@@@@@@@@@@@@@@@@@@      
                                                                -@@@@@@@@@@@@@@@@@@@@@@@@@@@@       
                                                             @@@@@@@@@@@@@@@@@@@@@@@@@@@@           
                                                         -@@@@@@@@@@@@@@@@@@@@@@@@@@@@@             
                                                      @@@@@@@@@@@@@@@@@@@@@@@@%  @@@@@@:            
                                                  -@@@@@@@@@@@@@@@@@@@@@@@@@      @@@@@@            
                                               @@@@@@@@@@@@@@@@@@@@@@@@@@@@@      @@@@@@            
                                           :@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@%     .@@@@@            
                                        @@@@@@@@@@@@@@@@@@.    :@@@@@@@@@@@@@      @@@@@            
                                    -@@@@@@@@@@@@@@@@.          @@@@@@@@@@@@@      @@@@@            
                                 @@@@@@@@@@@@@@@-.==:=         :@@@@@@@@@@@@@      @@@@@            
                             :@@@@@@@@@@@@@.=:======  =        @@@@@@@@@@@@@@      @@@@@            
                          @@@@@@@@@@@%      :==  -=-=..        @@@@@@@@@@@@@.      @@@@@            
                      :@@@@@@@@@#         = =-   =-= =.       @@@@@@@@@@@@@@       @@@@*            
                   @@@@@@@@@@@@@           -===- :..=.       @@@@@@@@@@@@@@.      *@@@@             
               :@@@@@@@@@@@@@@@@@            =: == =       .@@@@@@@@@@@@@@@       @@@@+             
            @@@@@- @@@@@@@@@@@@@@@@                       @@@@@@@@@@@@@@@@       #@@@@              
        :@@@.       @@@@@@@@@@@@@@@@@=                 %@@@@@@@@@@@@@@@@@        @@@@               
     @@.             @@@@@@@@@@@@@@@@@@@@%         @@@@@@@@@@@@@@@@@@@@@        @@@@                
 :.                   %@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@=        @@@@                 
                        @@@@@@@@@@@**@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@         @@@@                  
                          @@@@@    . @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@         :@@@@                   
                              .=====  @@@@@@@@@@@@@@@@@@@@@@@@@@@          @@@@                     
                    .       =========  @@@@@@@@@@@@@@@@@@@@@@@           *@@@-                      
                 -===   .== :========= =@@@@@@@@@@@@@@@@@@+            +@@@-                        
                  ==- =====- =========: #@@@@@@@@@@@+.               @@@@                           
                      :=====  =====                               #@@@-                             
                       ==.   :=.                              .@@@%                                 
                          ======                         -@@@@.                                     
                           ======                 ...                                               
                            ==                                                                      
                                                                                                    

░█████████                                                       ░██░██     ░██                       ░██       
░██     ░██                                                         ░██     ░██                       ░██       
░██     ░██  ░██████   ░██░████  ░██████   ░████████   ░███████  ░██░██     ░██  ░██████    ░███████  ░██    ░██
░█████████        ░██  ░███           ░██  ░██    ░██ ░██    ░██ ░██░██████████       ░██  ░██    ░██ ░██   ░██ 
░██          ░███████  ░██       ░███████  ░██    ░██ ░██    ░██ ░██░██     ░██  ░███████  ░██        ░███████  
░██         ░██   ░██  ░██      ░██   ░██  ░██    ░██ ░██    ░██ ░██░██     ░██ ░██   ░██  ░██    ░██ ░██   ░██ 
░██          ░█████░██ ░██       ░█████░██ ░██    ░██  ░███████  ░██░██     ░██  ░█████░██  ░███████  ░██    ░██
	`)
}

// DisplayCredit prints the open-source credit and contact information for ParanoiHack to w.
func DisplayCredit(w io.Writer) {
	fmt.Fprintln(w, `
		Open Source Software created and maintained by ParanoiHack
				https://paranoihack.ch
				sg@paranoihack.com
	`)
}

// DisplayFindings renders a formatted table of scan findings to w.
// Each row contains the engine name, severity, finding name, CWE, description,
// sink file path, sink line number, and remediation recommendation.
func DisplayFindings(w io.Writer, findings []models.Finding) {
	t := table.NewWriter()

	for _, finding := range findings {
		t.AppendRow(table.Row{
			finding.Engine,
			finding.Severity,
			finding.Name,
			finding.Cwe,
			finding.Description,
			finding.SinkFile,
			finding.SinkLine,
			finding.Recommendation,
			finding.Status,
		})
	}

	t.AppendHeader(table.Row{rowEngine, rowSeverity, rowName, rowCwe, rowDescription, rowSinkFile, rowSinkLine, rowRecommendation, rowStatus})

	t.SetCaption(fmt.Sprintf(caption, environment_variable.EnvironmentVariable["SCAN_DIR"]))
	t.SetColumnConfigs([]table.ColumnConfig{
		{Name: rowEngine, Align: text.AlignCenter, VAlign: text.VAlignMiddle, WidthMax: 10},
		{Name: rowSeverity, Align: text.AlignCenter, VAlign: text.VAlignMiddle, WidthMax: 10},
		{Name: rowCwe, Align: text.AlignCenter, VAlign: text.VAlignMiddle, WidthMax: 10},
		{Name: rowSinkFile, Align: text.AlignCenter, VAlign: text.VAlignMiddle, WidthMax: 50},
		{Name: rowSinkLine, Align: text.AlignCenter, VAlign: text.VAlignMiddle, WidthMax: 10},
		{Name: rowDescription, Align: text.AlignDefault, VAlign: text.VAlignMiddle, WidthMax: 50},
		{Name: rowRecommendation, Align: text.AlignDefault, VAlign: text.VAlignMiddle, WidthMax: 50},
		{Name: rowStatus, Align: text.AlignCenter, VAlign: text.VAlignMiddle, WidthMax: 10},
	})

	t.SetStyle(table.StyleLight)
	t.Style().Options.DrawBorder = true
	t.Style().Options.SeparateRows = true

	fmt.Fprintln(w, t.Render())
}

// DumpFindings writes the findings to w in the specified format.
// Supported formats are "json", "csv", and "raw" (table).
// An error is returned if the format is unsupported or if writing fails.
func DumpFindings(w io.Writer, findings []models.Finding, format string) error {
	switch format {
	case "json":
		return dumpFindingsJSON(w, findings)
	case "csv":
		return dumpFindingsCSV(w, findings)
	case "raw":
		dumpFindingsRaw(w, findings)
		return nil
	default:
		return fmt.Errorf(errUnsupportedFormat, format)
	}
}

func dumpFindingsJSON(w io.Writer, findings []models.Finding) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(findings)
}

func dumpFindingsCSV(w io.Writer, findings []models.Finding) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{rowEngine, rowSeverity, rowName, rowCwe, rowDescription, csvColSinkFile, csvColSinkLine, rowRecommendation, rowStatus}); err != nil {
		return err
	}
	for _, f := range findings {
		if err := cw.Write([]string{
			f.Engine,
			f.Severity,
			f.Name,
			f.Cwe,
			f.Description,
			f.SinkFile,
			strconv.Itoa(f.SinkLine),
			f.Recommendation,
			f.Status,
		}); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

func dumpFindingsRaw(w io.Writer, findings []models.Finding) {
	DisplayFindings(w, findings)
}
