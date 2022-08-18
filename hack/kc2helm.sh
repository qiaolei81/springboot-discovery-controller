#!/usr/bin/env bash

SRC_FOLDER="../config"
DST_FOLDER="../helm"

for f in $SRC_FOLDER/rbac/*.yaml; do if [[ $f != *kustomization.yaml* ]]; then echo -e "---"; cat $f; echo -e "\n"; fi; done | sed "s/namespace: system/namespace: {{ .Values.namespace }}/g" > $DST_FOLDER/templates/rbac.yaml
for f in $SRC_FOLDER/crd/bases/*.yaml; do if [[ $f != *kustomization.yaml* ]]; then echo -e "---"; cat $f; echo -e "\n"; fi; done | sed "s/namespace: system/namespace: {{ .Values.namespace }}/g" > $DST_FOLDER/templates/crd.yaml