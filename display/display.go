package display

import (
	"fmt"
	"scope-guardian/domains/models"
	environment_variable "scope-guardian/environnement_variable"

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
)

func DisplayBanner() {
	fmt.Println(`                                                                                   
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

func DisplayCredit() {
	fmt.Println(`
		Open Source Software created and maintained by ParanoiHack
				https://paranoihack.ch
				contact@paranoihack.com
	`)
}

func DisplayFindings(findings []models.Finding) {
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
		})
	}

	t.AppendHeader(table.Row{rowEngine, rowSeverity, rowName, rowCwe, rowDescription, rowSinkFile, rowSinkLine, rowRecommendation})

	t.SetCaption(fmt.Sprintf(caption, environment_variable.EnvironmentVariable["SCAN_DIR"]))
	t.SetColumnConfigs([]table.ColumnConfig{
		{Name: rowEngine, Align: text.AlignCenter, VAlign: text.VAlignMiddle, WidthMax: 10},
		{Name: rowSeverity, Align: text.AlignCenter, VAlign: text.VAlignMiddle, WidthMax: 10},
		{Name: rowCwe, Align: text.AlignCenter, VAlign: text.VAlignMiddle, WidthMax: 10},
		{Name: rowSinkFile, Align: text.AlignCenter, VAlign: text.VAlignMiddle, WidthMax: 50},
		{Name: rowSinkLine, Align: text.AlignCenter, VAlign: text.VAlignMiddle, WidthMax: 10},
		{Name: rowDescription, Align: text.AlignDefault, VAlign: text.VAlignMiddle, WidthMax: 50},
		{Name: rowRecommendation, Align: text.AlignDefault, VAlign: text.VAlignMiddle, WidthMax: 50},
	})

	t.SetStyle(table.StyleLight)
	t.Style().Options.DrawBorder = true
	t.Style().Options.SeparateRows = true

	fmt.Println(t.Render())
}
