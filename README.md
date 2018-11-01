# programming-k8s
a kubernetes programming example (via kubebuilder)

generate code via kubebuilder (refer: https://book.kubebuilder.io/quick_start.html)
```bash
$ ~/go/bin/kubebuilder_dir/bin/kubebuilder init --domain example.org --license apache2 --owner "interma"
$ ~/go/bin/kubebuilder_dir/bin/kubebuilder create api --group stats --version v1alpha1 --kind PodStats
```

