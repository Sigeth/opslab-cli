package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	"github.com/Sigeth/opslab-core-go/api/v1alpha1"
	"github.com/Sigeth/opslab-core-go/pkg/challenge"
	"github.com/Sigeth/opslab-core-go/pkg/project"
)

// defaultSystemNamespace is where the platform (and its operator) is installed.
const defaultSystemNamespace = project.Name + "-system"

func newChallengeManifestCommand() *cobra.Command {
	namespace := defaultSystemNamespace
	cmd := &cobra.Command{
		Use:   "manifest [DIR]",
		Short: "Print the Kubernetes objects installing a local challenge (development)",
		Long: "Prints a ConfigMap holding the package files and a Challenge embedding the manifest,\n" +
			"for clusters without a registry:\n\n" +
			"  " + project.Name + " challenge manifest . | kubectl apply -f -\n\n" +
			"Inline challenges are not signed: use them for development only.",
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := "."
			if len(args) == 1 {
				dir = args[0]
			}
			return manifest(cmd.OutOrStdout(), dir, namespace)
		},
	}
	cmd.Flags().StringVarP(&namespace, "namespace", "n", namespace, "namespace of the operator, which holds the package files")
	return cmd
}

func manifest(out io.Writer, dir, namespace string) error {
	files := os.DirFS(dir)
	pkg, err := challenge.Load(files)
	if err != nil {
		return fmt.Errorf("%s: %w", dir, err)
	}
	if errs := challenge.Validate(pkg, files); len(errs) > 0 {
		return fmt.Errorf("%w: %s: %v (run '%s challenge lint')", errInvalid, dir, errs.ToAggregate(), project.Name)
	}
	bundle, err := challenge.Bundle(pkg, files)
	if err != nil {
		return err
	}
	delete(bundle, challenge.ManifestFile) // embedded in the Challenge
	data, binaryData, err := challenge.ConfigMapData(bundle)
	if err != nil {
		return err
	}

	labels := map[string]string{project.LabelPrefix + "challenge": pkg.Metadata.Name}
	cm := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{APIVersion: "v1", Kind: "ConfigMap"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "challenge-" + pkg.Metadata.Name + "-files",
			Namespace: namespace,
			Labels:    labels,
		},
		Data:       data,
		BinaryData: binaryData,
	}
	ch := &v1alpha1.Challenge{
		TypeMeta:   metav1.TypeMeta{APIVersion: v1alpha1.GroupVersion.String(), Kind: "Challenge"},
		ObjectMeta: metav1.ObjectMeta{Name: pkg.Metadata.Name, Labels: labels},
		Spec: v1alpha1.ChallengeSpec{Source: v1alpha1.ChallengeSource{
			Inline: &v1alpha1.InlinePackage{Package: *pkg, FilesConfigMap: cm.Name},
		}},
	}

	fmt.Fprintf(out, "# %s %s from %s (inline, unsigned: development only)\n", pkg.Metadata.Name, pkg.Metadata.Version, dir)
	for i, obj := range []any{cm, ch} {
		doc, err := yaml.Marshal(obj)
		if err != nil {
			return err
		}
		if i > 0 {
			fmt.Fprintln(out, "---")
		}
		if _, err := out.Write(doc); err != nil {
			return err
		}
	}
	return nil
}
