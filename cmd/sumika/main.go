package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/mskasa/sumika/internal/config"
	"github.com/mskasa/sumika/internal/git"
	"github.com/mskasa/sumika/internal/launcher"
	"github.com/mskasa/sumika/internal/project"
	"github.com/mskasa/sumika/internal/server"
	"github.com/spf13/cobra"
)

func main() {
	if err := newRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "sumika",
		Short: "Personal project hub for solo developers",
	}
	root.AddCommand(
		newInitCmd(),
		newAddCmd(),
		newListCmd(),
		newOpenCmd(),
		newStatusCmd(),
		newRemoveCmd(),
		newServeCmd(),
	)
	return root
}

func newInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize sumika config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := config.DefaultPath()
			if err != nil {
				return err
			}
			if _, err := os.Stat(path); err == nil {
				return fmt.Errorf("config already exists at %s", path)
			}
			cfg := &config.Config{
				Version: 1,
				Settings: config.Settings{
					Port: 8964,
				},
			}
			if err := cfg.SaveTo(path); err != nil {
				return err
			}
			fmt.Printf("Initialized sumika config at %s\n", path)
			return nil
		},
	}
}

func newAddCmd() *cobra.Command {
	var name, description string
	cmd := &cobra.Command{
		Use:   "add <path>",
		Short: "Register a project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			info, err := os.Stat(path)
			if err != nil || !info.IsDir() {
				return fmt.Errorf("%q is not a valid directory", path)
			}
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			if err := project.Add(cfg, path, name, description); err != nil {
				return err
			}
			if err := cfg.Save(); err != nil {
				return err
			}
			resolvedName := name
			if resolvedName == "" {
				abs, _ := filepath.Abs(path)
				resolvedName = filepath.Base(abs)
			}
			fmt.Printf("Added project %q\n", resolvedName)
			return nil
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "Project name (default: directory name)")
	cmd.Flags().StringVar(&description, "description", "", "Project description")
	return cmd
}

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List registered projects",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			if len(cfg.Projects) == 0 {
				fmt.Println("No projects registered. Run `sumika add <path>` to register one.")
				return nil
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "NAME\tPATH\tDESCRIPTION\tTAGS")
			for _, p := range cfg.Projects {
				tags := strings.Join(p.Tags, ", ")
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", p.Name, p.Path, p.Description, tags)
			}
			return w.Flush()
		},
	}
}

func newRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove <name>",
		Short: "Unregister a project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			if err := project.Remove(cfg, args[0]); err != nil {
				return err
			}
			if err := cfg.Save(); err != nil {
				return err
			}
			fmt.Printf("Removed project %q\n", args[0])
			return nil
		},
	}
}

func newOpenCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "open <name>",
		Short: "Open project in editor and AI tool",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			p, err := project.Find(cfg, args[0])
			if err != nil {
				return err
			}
			editor := ""
			if p.Launch.Editor {
				editor = cfg.Settings.Editor
			}
			aiTool := ""
			if p.Launch.AI {
				aiTool = cfg.Settings.AITool
			}
			slog.Info("opening project", "name", p.Name, "path", p.Path)
			return launcher.Open(p.Path, editor, aiTool, p.Launch.Commands)
		},
	}
}

func newStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show git status of all projects",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			if len(cfg.Projects) == 0 {
				fmt.Println("No projects registered.")
				return nil
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "NAME\tLAST COMMIT\tCHANGES")
			for _, p := range cfg.Projects {
				st, err := git.GetStatus(p.Path, nil)
				if err != nil {
					fmt.Fprintf(w, "%s\t-\terror: %v\n", p.Name, err)
					continue
				}
				if !st.IsRepo {
					fmt.Fprintf(w, "%s\t-\tnot a git repository\n", p.Name)
					continue
				}
				changes := "clean"
				if st.UncommittedCount > 0 {
					changes = fmt.Sprintf("%d uncommitted", st.UncommittedCount)
				}
				fmt.Fprintf(w, "%s\t%s\t%s\n", p.Name, st.LastCommit, changes)
			}
			return w.Flush()
		},
	}
}

func newServeCmd() *cobra.Command {
	var port int
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the web dashboard",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			if port == 0 {
				port = cfg.Settings.Port
			}
			s := server.New(cfg)
			return s.Run(port)
		},
	}
	cmd.Flags().IntVar(&port, "port", 0, "Port to listen on (default: config value or 8964)")
	return cmd
}
