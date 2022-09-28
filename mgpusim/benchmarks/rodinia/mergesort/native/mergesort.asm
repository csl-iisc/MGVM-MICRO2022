	.text
	.hsa_code_object_version 2,1
	.hsa_code_object_isa 8,0,3,"AMD","AMDGPU"
	.protected	mergeSortFirst  ; -- Begin function mergeSortFirst
	.globl	mergeSortFirst
	.p2align	8
	.type	mergeSortFirst,@function
	.amdgpu_hsa_kernel mergeSortFirst
mergeSortFirst:                         ; @mergeSortFirst
mergeSortFirst$local:
	.amd_kernel_code_t
		amd_code_version_major = 1
		amd_code_version_minor = 2
		amd_machine_kind = 1
		amd_machine_version_major = 8
		amd_machine_version_minor = 0
		amd_machine_version_stepping = 3
		kernel_code_entry_byte_offset = 256
		kernel_code_prefetch_byte_size = 0
		granulated_workitem_vgpr_count = 2
		granulated_wavefront_sgpr_count = 1
		priority = 0
		float_mode = 192
		priv = 0
		enable_dx10_clamp = 1
		debug_mode = 0
		enable_ieee_mode = 1
		enable_wgp_mode = 0
		enable_mem_ordered = 0
		enable_fwd_progress = 0
		enable_sgpr_private_segment_wave_byte_offset = 0
		user_sgpr_count = 8
		enable_trap_handler = 0
		enable_sgpr_workgroup_id_x = 1
		enable_sgpr_workgroup_id_y = 0
		enable_sgpr_workgroup_id_z = 0
		enable_sgpr_workgroup_info = 0
		enable_vgpr_workitem_id = 0
		enable_exception_msb = 0
		granulated_lds_size = 0
		enable_exception = 0
		enable_sgpr_private_segment_buffer = 1
		enable_sgpr_dispatch_ptr = 1
		enable_sgpr_queue_ptr = 0
		enable_sgpr_kernarg_segment_ptr = 1
		enable_sgpr_dispatch_id = 0
		enable_sgpr_flat_scratch_init = 0
		enable_sgpr_private_segment_size = 0
		enable_sgpr_grid_workgroup_count_x = 0
		enable_sgpr_grid_workgroup_count_y = 0
		enable_sgpr_grid_workgroup_count_z = 0
		enable_wavefront_size32 = 0
		enable_ordered_append_gds = 0
		private_element_size = 1
		is_ptr64 = 1
		is_dynamic_callstack = 0
		is_debug_enabled = 0
		is_xnack_enabled = 0
		workitem_private_segment_byte_size = 0
		workgroup_group_segment_byte_size = 0
		gds_segment_byte_size = 0
		kernarg_segment_byte_size = 80
		workgroup_fbarrier_count = 0
		wavefront_sgpr_count = 11
		workitem_vgpr_count = 10
		reserved_vgpr_first = 0
		reserved_vgpr_count = 0
		reserved_sgpr_first = 0
		reserved_sgpr_count = 0
		debug_wavefront_private_segment_offset_sgpr = 0
		debug_private_segment_buffer_sgpr = 0
		kernarg_segment_alignment = 4
		group_segment_alignment = 4
		private_segment_alignment = 4
		wavefront_size = 6
		call_convention = -1
		runtime_loader_kernel_symbol = 0
	.end_amd_kernel_code_t
; %bb.0:
	s_load_dword s2, s[6:7], 0x10
	s_load_dword s0, s[4:5], 0x4
	v_mov_b32_e32 v1, 0
	v_mov_b32_e32 v2, s8
	s_waitcnt lgkmcnt(0)
	s_and_b32 s0, s0, 0xffff
	v_mad_i64_i32 v[0:1], s[0:1], s0, v2, v[0:1]
	s_ashr_i32 s0, s2, 31
	s_lshr_b32 s0, s0, 30
	s_add_i32 s2, s2, s0
	s_ashr_i32 s0, s2, 2
	s_ashr_i32 s1, s0, 31
	v_cmp_gt_u64_e32 vcc, s[0:1], v[0:1]
	s_and_saveexec_b64 s[0:1], vcc
	s_cbranch_execz BB0_2
; %bb.1:
	s_load_dwordx4 s[0:3], s[6:7], 0x0
	v_lshlrev_b64 v[4:5], 4, v[0:1]
	s_waitcnt lgkmcnt(0)
	v_mov_b32_e32 v1, s1
	v_add_u32_e32 v0, vcc, s0, v4
	v_addc_u32_e32 v1, vcc, v1, v5, vcc
	flat_load_dwordx4 v[0:3], v[0:1]
	v_mov_b32_e32 v6, s3
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_cmp_gt_f32_e32 vcc, v0, v1
	v_cndmask_b32_e32 v7, v0, v1, vcc
	v_cmp_gt_f32_e32 vcc, v1, v0
	v_cndmask_b32_e32 v1, v0, v1, vcc
	v_cmp_gt_f32_e32 vcc, v2, v3
	v_cndmask_b32_e32 v8, v2, v3, vcc
	v_cmp_gt_f32_e32 vcc, v3, v2
	v_cndmask_b32_e32 v2, v2, v3, vcc
	v_cmp_gt_f32_e32 vcc, v7, v8
	v_cndmask_b32_e32 v0, v7, v8, vcc
	v_cmp_gt_f32_e32 vcc, v1, v2
	v_cndmask_b32_e32 v9, v1, v2, vcc
	v_cmp_gt_f32_e32 vcc, v8, v7
	v_cndmask_b32_e32 v7, v7, v8, vcc
	v_cmp_gt_f32_e32 vcc, v2, v1
	v_cndmask_b32_e32 v3, v1, v2, vcc
	v_cmp_gt_f32_e32 vcc, v9, v7
	v_cndmask_b32_e32 v1, v9, v7, vcc
	v_cmp_gt_f32_e32 vcc, v7, v9
	v_cndmask_b32_e32 v2, v9, v7, vcc
	v_add_u32_e32 v4, vcc, s2, v4
	v_addc_u32_e32 v5, vcc, v6, v5, vcc
	flat_store_dwordx4 v[4:5], v[0:3]
BB0_2:
	s_endpgm
.Lfunc_end0:
	.size	mergeSortFirst, .Lfunc_end0-mergeSortFirst
                                        ; -- End function
	.section	.AMDGPU.csdata
; Kernel info:
; codeLenInByte = 224
; NumSgprs: 11
; NumVgprs: 10
; ScratchSize: 0
; MemoryBound: 0
; FloatMode: 192
; IeeeMode: 1
; LDSByteSize: 0 bytes/workgroup (compile time only)
; SGPRBlocks: 1
; VGPRBlocks: 2
; NumSGPRsForWavesPerEU: 11
; NumVGPRsForWavesPerEU: 10
; Occupancy: 10
; WaveLimiterHint : 1
; COMPUTE_PGM_RSRC2:USER_SGPR: 8
; COMPUTE_PGM_RSRC2:TRAP_HANDLER: 0
; COMPUTE_PGM_RSRC2:TGID_X_EN: 1
; COMPUTE_PGM_RSRC2:TGID_Y_EN: 0
; COMPUTE_PGM_RSRC2:TGID_Z_EN: 0
; COMPUTE_PGM_RSRC2:TIDIG_COMP_CNT: 0
	.text
	.protected	mergeSortPass   ; -- Begin function mergeSortPass
	.globl	mergeSortPass
	.p2align	8
	.type	mergeSortPass,@function
	.amdgpu_hsa_kernel mergeSortPass
mergeSortPass:                          ; @mergeSortPass
mergeSortPass$local:
	.amd_kernel_code_t
		amd_code_version_major = 1
		amd_code_version_minor = 2
		amd_machine_kind = 1
		amd_machine_version_major = 8
		amd_machine_version_minor = 0
		amd_machine_version_stepping = 3
		kernel_code_entry_byte_offset = 256
		kernel_code_prefetch_byte_size = 0
		granulated_workitem_vgpr_count = 9
		granulated_wavefront_sgpr_count = 4
		priority = 0
		float_mode = 192
		priv = 0
		enable_dx10_clamp = 1
		debug_mode = 0
		enable_ieee_mode = 1
		enable_wgp_mode = 0
		enable_mem_ordered = 0
		enable_fwd_progress = 0
		enable_sgpr_private_segment_wave_byte_offset = 0
		user_sgpr_count = 8
		enable_trap_handler = 0
		enable_sgpr_workgroup_id_x = 1
		enable_sgpr_workgroup_id_y = 0
		enable_sgpr_workgroup_id_z = 0
		enable_sgpr_workgroup_info = 0
		enable_vgpr_workitem_id = 0
		enable_exception_msb = 0
		granulated_lds_size = 0
		enable_exception = 0
		enable_sgpr_private_segment_buffer = 1
		enable_sgpr_dispatch_ptr = 1
		enable_sgpr_queue_ptr = 0
		enable_sgpr_kernarg_segment_ptr = 1
		enable_sgpr_dispatch_id = 0
		enable_sgpr_flat_scratch_init = 0
		enable_sgpr_private_segment_size = 0
		enable_sgpr_grid_workgroup_count_x = 0
		enable_sgpr_grid_workgroup_count_y = 0
		enable_sgpr_grid_workgroup_count_z = 0
		enable_wavefront_size32 = 0
		enable_ordered_append_gds = 0
		private_element_size = 1
		is_ptr64 = 1
		is_dynamic_callstack = 0
		is_debug_enabled = 0
		is_xnack_enabled = 0
		workitem_private_segment_byte_size = 0
		workgroup_group_segment_byte_size = 0
		gds_segment_byte_size = 0
		kernarg_segment_byte_size = 88
		workgroup_fbarrier_count = 0
		wavefront_sgpr_count = 40
		workitem_vgpr_count = 37
		reserved_vgpr_first = 0
		reserved_vgpr_count = 0
		reserved_sgpr_first = 0
		reserved_sgpr_count = 0
		debug_wavefront_private_segment_offset_sgpr = 0
		debug_private_segment_buffer_sgpr = 0
		kernarg_segment_alignment = 4
		group_segment_alignment = 4
		private_segment_alignment = 4
		wavefront_size = 6
		call_convention = -1
		runtime_loader_kernel_symbol = 0
	.end_amd_kernel_code_t
; %bb.0:
	s_load_dwordx2 s[10:11], s[6:7], 0x10
	s_load_dword s0, s[6:7], 0x20
	s_waitcnt lgkmcnt(0)
	s_ashr_i32 s2, s11, 31
	s_add_i32 s1, s11, s2
	s_xor_b32 s3, s1, s2
	v_cvt_f32_u32_e32 v1, s3
	s_load_dword s1, s[4:5], 0x4
	v_rcp_iflag_f32_e32 v1, v1
	s_waitcnt lgkmcnt(0)
	s_and_b32 s1, s1, 0xffff
	s_mul_i32 s8, s8, s1
	v_mul_f32_e32 v1, 0x4f800000, v1
	v_cvt_u32_f32_e32 v1, v1
	s_add_i32 s0, s0, s8
	v_add_u32_e32 v0, vcc, s0, v0
	v_ashrrev_i32_e32 v4, 31, v0
	v_mul_lo_u32 v2, v1, s3
	v_mul_hi_u32 v3, v1, s3
	v_sub_u32_e32 v5, vcc, 0, v2
	v_cmp_eq_u32_e64 s[0:1], 0, v3
	v_cndmask_b32_e64 v2, v2, v5, s[0:1]
	v_mul_hi_u32 v2, v2, v1
	v_add_u32_e32 v3, vcc, v0, v4
	v_xor_b32_e32 v3, v3, v4
	v_add_u32_e32 v5, vcc, v2, v1
	v_subrev_u32_e32 v1, vcc, v2, v1
	v_cndmask_b32_e64 v1, v1, v5, s[0:1]
	v_mul_hi_u32 v1, v1, v3
	v_xor_b32_e32 v2, s2, v4
	v_mul_lo_u32 v4, v1, s3
	v_add_u32_e32 v5, vcc, -1, v1
	v_sub_u32_e32 v6, vcc, v3, v4
	v_cmp_ge_u32_e32 vcc, v3, v4
	v_cmp_le_u32_e64 s[0:1], s3, v6
	v_add_u32_e64 v3, s[2:3], 1, v1
	s_and_b64 s[0:1], s[0:1], vcc
	v_cndmask_b32_e64 v1, v1, v3, s[0:1]
	v_cndmask_b32_e32 v1, v5, v1, vcc
	v_xor_b32_e32 v1, v1, v2
	v_sub_u32_e32 v1, vcc, v1, v2
	s_movk_i32 s0, 0x400
	v_cmp_gt_i32_e32 vcc, s0, v1
	s_and_saveexec_b64 s[0:1], vcc
	s_cbranch_execz BB1_24
; %bb.1:
	v_mul_lo_u32 v3, v1, s11
	s_load_dwordx2 s[0:1], s[6:7], 0x18
	v_ashrrev_i32_e32 v2, 31, v1
	v_sub_u32_e32 v5, vcc, v0, v3
	v_lshlrev_b64 v[0:1], 2, v[1:2]
	s_waitcnt lgkmcnt(0)
	v_mov_b32_e32 v2, s1
	v_add_u32_e32 v3, vcc, s0, v0
	v_addc_u32_e32 v4, vcc, v2, v1, vcc
	v_mul_lo_u32 v0, v5, s10
	flat_load_dwordx2 v[5:6], v[3:4]
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_add_u32_e32 v9, vcc, v5, v0
	v_cmp_lt_i32_e32 vcc, v9, v6
	s_and_b64 exec, exec, vcc
	s_cbranch_execz BB1_24
; %bb.2:
	s_load_dwordx4 s[4:7], s[6:7], 0x0
	s_lshr_b32 s0, s10, 31
	s_add_i32 s0, s10, s0
	v_ashrrev_i32_e32 v10, 31, v9
	s_ashr_i32 s8, s0, 1
	v_lshlrev_b64 v[0:1], 4, v[9:10]
	v_add_u32_e32 v2, vcc, s8, v9
	s_waitcnt lgkmcnt(0)
	v_mov_b32_e32 v5, s7
	v_add_u32_e32 v0, vcc, s6, v0
	v_addc_u32_e32 v1, vcc, v5, v1, vcc
	v_add_u32_e32 v3, vcc, 4, v3
	v_addc_u32_e32 v4, vcc, 0, v4, vcc
	v_cmp_ge_i32_e32 vcc, v2, v6
	s_and_saveexec_b64 s[0:1], vcc
	s_xor_b64 s[0:1], exec, s[0:1]
	s_cbranch_execz BB1_6
; %bb.3:                                ; %.preheader.preheader
	v_mov_b32_e32 v6, v1
	s_mov_b32 s6, 0
	s_mov_b64 s[2:3], 0
	v_mov_b32_e32 v5, v0
BB1_4:                                  ; %.preheader
                                        ; =>This Inner Loop Header: Depth=1
	v_add_u32_e32 v7, vcc, s6, v9
	v_ashrrev_i32_e32 v8, 31, v7
	v_lshlrev_b64 v[7:8], 4, v[7:8]
	v_mov_b32_e32 v11, s5
	v_add_u32_e32 v7, vcc, s4, v7
	v_addc_u32_e32 v8, vcc, v11, v8, vcc
	flat_load_dwordx4 v[11:14], v[7:8]
	s_add_i32 s6, s6, 1
	s_waitcnt vmcnt(0) lgkmcnt(0)
	flat_store_dwordx4 v[5:6], v[11:14]
	flat_load_dword v7, v[3:4]
	v_add_u32_e32 v5, vcc, 16, v5
	v_addc_u32_e32 v6, vcc, 0, v6, vcc
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_sub_u32_e32 v7, vcc, v7, v9
	v_cmp_ge_i32_e32 vcc, s6, v7
	s_or_b64 s[2:3], vcc, s[2:3]
	s_andn2_b64 exec, exec, s[2:3]
	s_cbranch_execnz BB1_4
; %bb.5:                                ; %Flow
	s_or_b64 exec, exec, s[2:3]
BB1_6:                                  ; %Flow271
	s_or_saveexec_b64 s[0:1], s[0:1]
	s_xor_b64 exec, exec, s[0:1]
	s_cbranch_execz BB1_24
; %bb.7:
	v_lshlrev_b64 v[5:6], 4, v[9:10]
	s_ashr_i32 s9, s8, 31
	v_mov_b32_e32 v7, s5
	v_add_u32_e32 v5, vcc, s4, v5
	s_lshl_b64 s[0:1], s[8:9], 4
	v_addc_u32_e32 v6, vcc, v7, v6, vcc
	v_mov_b32_e32 v8, s1
	v_add_u32_e32 v7, vcc, s0, v5
	v_addc_u32_e32 v8, vcc, v6, v8, vcc
	flat_load_dwordx4 v[13:16], v[5:6]
	flat_load_dwordx4 v[5:8], v[7:8]
	v_mov_b32_e32 v17, 0
	v_add_u32_e32 v9, vcc, 1, v9
	s_mov_b64 s[6:7], 0
	v_mov_b32_e32 v12, v17
	v_mov_b32_e32 v10, v17
                                        ; implicit-def: $sgpr10_sgpr11
	s_branch BB1_9
BB1_8:                                  ; %Flow270
                                        ;   in Loop: Header=BB1_9 Depth=1
	s_or_b64 exec, exec, s[2:3]
	s_and_b64 s[0:1], exec, s[0:1]
	s_or_b64 s[6:7], s[0:1], s[6:7]
	v_mov_b32_e32 v12, v11
	s_andn2_b64 exec, exec, s[6:7]
	s_cbranch_execz BB1_23
BB1_9:                                  ; =>This Loop Header: Depth=1
                                        ;     Child Loop BB1_11 Depth 2
	v_add_u32_e32 v11, vcc, 1, v12
	v_add_u32_e32 v19, vcc, v11, v2
	v_ashrrev_i32_e32 v20, 31, v19
	v_lshlrev_b64 v[20:21], 4, v[19:20]
	v_mov_b32_e32 v18, s5
	v_add_u32_e32 v20, vcc, s4, v20
	v_addc_u32_e32 v21, vcc, v18, v21, vcc
	v_mov_b32_e32 v18, 0
	v_lshlrev_b64 v[22:23], 4, v[17:18]
	s_waitcnt vmcnt(1) lgkmcnt(1)
	v_mov_b32_e32 v31, v16
	v_add_u32_e64 v22, s[0:1], v0, v22
	v_addc_u32_e64 v23, s[0:1], v1, v23, s[0:1]
	v_cmp_le_i32_e32 vcc, s8, v11
	v_add_u32_e64 v32, s[0:1], 1, v17
	s_mov_b64 s[16:17], 0
	s_mov_b64 s[22:23], s[10:11]
	v_mov_b32_e32 v30, v15
	v_mov_b32_e32 v29, v14
	v_mov_b32_e32 v28, v13
	v_mov_b32_e32 v18, v10
                                        ; implicit-def: $sgpr14_sgpr15
                                        ; implicit-def: $sgpr12_sgpr13
                                        ; implicit-def: $sgpr18_sgpr19
                                        ; implicit-def: $sgpr20_sgpr21
	s_branch BB1_11
BB1_10:                                 ; %Flow267
                                        ;   in Loop: Header=BB1_11 Depth=2
	s_or_b64 exec, exec, s[26:27]
	v_cmp_ge_f32_e64 s[0:1], v31, v5
	v_cndmask_b32_e64 v5, v5, v31, s[0:1]
	v_cmp_ge_f32_e64 s[0:1], v30, v6
	v_cndmask_b32_e64 v6, v6, v30, s[0:1]
	v_cmp_ge_f32_e64 s[0:1], v29, v7
	v_cndmask_b32_e64 v7, v7, v29, s[0:1]
	v_cmp_ge_f32_e64 s[0:1], v28, v8
	v_cndmask_b32_e64 v8, v8, v28, s[0:1]
	v_cmp_gt_f32_e64 s[0:1], v5, v6
	v_cndmask_b32_e64 v28, v5, v6, s[0:1]
	v_cmp_gt_f32_e64 s[0:1], v6, v5
	v_cndmask_b32_e64 v6, v5, v6, s[0:1]
	v_cmp_gt_f32_e64 s[0:1], v7, v8
	v_cndmask_b32_e64 v29, v7, v8, s[0:1]
	v_cmp_gt_f32_e64 s[0:1], v8, v7
	v_cndmask_b32_e64 v7, v7, v8, s[0:1]
	v_cmp_gt_f32_e64 s[0:1], v28, v29
	v_cndmask_b32_e64 v5, v28, v29, s[0:1]
	v_cmp_gt_f32_e64 s[0:1], v6, v7
	v_cndmask_b32_e64 v30, v6, v7, s[0:1]
	v_cmp_gt_f32_e64 s[0:1], v29, v28
	v_cndmask_b32_e64 v28, v28, v29, s[0:1]
	v_cmp_gt_f32_e64 s[0:1], v7, v6
	v_cndmask_b32_e64 v8, v6, v7, s[0:1]
	v_cmp_gt_f32_e64 s[0:1], v30, v28
	v_cndmask_b32_e64 v6, v30, v28, s[0:1]
	v_cmp_gt_f32_e64 s[0:1], v28, v30
	v_cndmask_b32_e64 v7, v30, v28, s[0:1]
	s_and_b64 s[0:1], exec, s[2:3]
	s_or_b64 s[16:17], s[0:1], s[16:17]
	s_andn2_b64 s[0:1], s[12:13], exec
	s_and_b64 s[2:3], s[18:19], exec
	s_or_b64 s[12:13], s[0:1], s[2:3]
	s_andn2_b64 s[0:1], s[10:11], exec
	s_and_b64 s[2:3], s[22:23], exec
	s_waitcnt vmcnt(2) lgkmcnt(2)
	v_mov_b32_e32 v31, v27
	s_or_b64 s[10:11], s[0:1], s[2:3]
	s_andn2_b64 s[0:1], s[14:15], exec
	s_and_b64 s[2:3], s[20:21], exec
	s_or_b64 s[14:15], s[0:1], s[2:3]
	v_mov_b32_e32 v30, v26
	v_mov_b32_e32 v29, v25
	v_mov_b32_e32 v28, v24
	s_andn2_b64 exec, exec, s[16:17]
	s_cbranch_execz BB1_19
BB1_11:                                 ;   Parent Loop BB1_9 Depth=1
                                        ; =>  This Inner Loop Header: Depth=2
	v_mov_b32_e32 v10, v18
	s_waitcnt vmcnt(1) lgkmcnt(1)
	v_add_u32_e64 v13, s[0:1], v9, v10
	v_ashrrev_i32_e32 v14, 31, v13
	v_lshlrev_b64 v[13:14], 4, v[13:14]
	v_mov_b32_e32 v15, s5
	v_add_u32_e64 v13, s[0:1], s4, v13
	v_addc_u32_e64 v14, s[0:1], v15, v14, s[0:1]
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_cmp_lt_f32_e64 s[0:1], v28, v8
	v_cndmask_b32_e64 v15, v8, v28, s[0:1]
	v_cmp_lt_f32_e64 s[0:1], v29, v7
	v_cndmask_b32_e64 v16, v7, v29, s[0:1]
	v_cmp_gt_f32_e64 s[0:1], v15, v16
	v_cndmask_b32_e64 v17, v15, v16, s[0:1]
	v_cmp_gt_f32_e64 s[0:1], v16, v15
	v_cndmask_b32_e64 v15, v15, v16, s[0:1]
	v_cmp_lt_f32_e64 s[0:1], v30, v6
	v_cndmask_b32_e64 v16, v6, v30, s[0:1]
	v_cmp_lt_f32_e64 s[0:1], v31, v5
	v_cndmask_b32_e64 v18, v5, v31, s[0:1]
	v_cmp_gt_f32_e64 s[0:1], v16, v18
	v_cndmask_b32_e64 v24, v16, v18, s[0:1]
	v_cmp_gt_f32_e64 s[0:1], v18, v16
	v_cndmask_b32_e64 v16, v16, v18, s[0:1]
	v_cmp_gt_f32_e64 s[0:1], v15, v16
	v_cndmask_b32_e64 v18, v15, v16, s[0:1]
	v_cmp_gt_f32_e64 s[0:1], v16, v15
	v_cndmask_b32_e64 v36, v15, v16, s[0:1]
	v_cmp_gt_f32_e64 s[0:1], v17, v24
	v_cndmask_b32_e64 v33, v17, v24, s[0:1]
	v_cmp_gt_f32_e64 s[0:1], v24, v17
	v_cndmask_b32_e64 v15, v17, v24, s[0:1]
	v_cmp_gt_f32_e64 s[0:1], v18, v15
	v_cndmask_b32_e64 v34, v18, v15, s[0:1]
	v_cmp_gt_f32_e64 s[0:1], v15, v18
	v_cndmask_b32_e64 v35, v18, v15, s[0:1]
	flat_load_dwordx4 v[24:27], v[13:14]
	flat_load_dwordx4 v[13:16], v[20:21]
	v_add_u32_e64 v18, s[0:1], 1, v10
	v_cmp_gt_i32_e64 s[0:1], s8, v18
	s_mov_b64 s[28:29], 0
	s_mov_b64 s[24:25], 0
	v_mov_b32_e32 v17, v32
	flat_store_dwordx4 v[22:23], v[33:36]
	s_and_saveexec_b64 s[2:3], vcc
	s_xor_b64 s[2:3], exec, s[2:3]
; %bb.12:                               ;   in Loop: Header=BB1_11 Depth=2
	s_and_b64 s[24:25], s[0:1], exec
; %bb.13:                               ; %Flow264
                                        ;   in Loop: Header=BB1_11 Depth=2
	s_or_saveexec_b64 s[26:27], s[2:3]
	s_andn2_b64 s[2:3], s[22:23], exec
	s_and_b64 s[22:23], s[22:23], exec
	s_or_b64 s[22:23], s[2:3], s[22:23]
	s_andn2_b64 s[20:21], s[20:21], exec
	s_xor_b64 exec, exec, s[26:27]
	s_cbranch_execz BB1_17
; %bb.14:                               ;   in Loop: Header=BB1_11 Depth=2
	flat_load_dword v32, v[3:4]
	s_mov_b64 s[30:31], -1
	s_mov_b64 s[28:29], 0
	s_mov_b64 s[34:35], s[24:25]
	s_waitcnt vmcnt(0) lgkmcnt(0)
	v_cmp_ge_i32_e64 s[2:3], v19, v32
	s_and_saveexec_b64 s[36:37], s[0:1]
; %bb.15:                               ;   in Loop: Header=BB1_11 Depth=2
	s_andn2_b64 s[0:1], s[24:25], exec
	s_and_b64 s[34:35], s[2:3], exec
	s_mov_b64 s[28:29], exec
	s_xor_b64 s[30:31], exec, -1
	s_or_b64 s[34:35], s[0:1], s[34:35]
; %bb.16:                               ; %Flow266
                                        ;   in Loop: Header=BB1_11 Depth=2
	s_or_b64 exec, exec, s[36:37]
	s_andn2_b64 s[0:1], s[22:23], exec
	s_and_b64 s[2:3], s[2:3], exec
	s_or_b64 s[22:23], s[0:1], s[2:3]
	s_andn2_b64 s[0:1], s[20:21], exec
	s_and_b64 s[2:3], s[30:31], exec
	s_or_b64 s[20:21], s[0:1], s[2:3]
	s_andn2_b64 s[0:1], s[24:25], exec
	s_and_b64 s[2:3], s[34:35], exec
	s_and_b64 s[28:29], s[28:29], exec
	s_or_b64 s[24:25], s[0:1], s[2:3]
BB1_17:                                 ; %Flow265
                                        ;   in Loop: Header=BB1_11 Depth=2
	s_or_b64 exec, exec, s[26:27]
	s_andn2_b64 s[0:1], s[18:19], exec
	s_and_b64 s[18:19], s[28:29], exec
	s_mov_b64 s[2:3], -1
	s_or_b64 s[18:19], s[0:1], s[18:19]
                                        ; implicit-def: $vgpr32
	s_and_saveexec_b64 s[26:27], s[24:25]
	s_cbranch_execz BB1_10
; %bb.18:                               ; %.backedge
                                        ;   in Loop: Header=BB1_11 Depth=2
	v_add_u32_e64 v22, s[0:1], 16, v22
	v_addc_u32_e64 v23, s[0:1], 0, v23, s[0:1]
	v_add_u32_e64 v32, s[0:1], 1, v17
	s_xor_b64 s[2:3], exec, -1
	s_andn2_b64 s[18:19], s[18:19], exec
	s_branch BB1_10
BB1_19:                                 ; %Flow268
                                        ;   in Loop: Header=BB1_9 Depth=1
	s_or_b64 exec, exec, s[16:17]
	s_mov_b64 s[0:1], -1
	s_and_saveexec_b64 s[2:3], s[14:15]
	s_xor_b64 s[2:3], exec, s[2:3]
; %bb.20:                               ;   in Loop: Header=BB1_9 Depth=1
	s_orn2_b64 s[0:1], s[10:11], exec
; %bb.21:                               ; %Flow269
                                        ;   in Loop: Header=BB1_9 Depth=1
	s_or_b64 exec, exec, s[2:3]
	s_and_saveexec_b64 s[2:3], s[12:13]
	s_cbranch_execz BB1_8
; %bb.22:                               ;   in Loop: Header=BB1_9 Depth=1
	s_waitcnt vmcnt(1) lgkmcnt(1)
	v_cmp_lt_f32_e32 vcc, v24, v13
	v_cndmask_b32_e32 v16, v16, v27, vcc
	v_cndmask_b32_e32 v15, v15, v26, vcc
	v_cndmask_b32_e32 v14, v14, v25, vcc
	v_cndmask_b32_e32 v13, v13, v24, vcc
	v_cndmask_b32_e32 v11, v11, v12, vcc
	v_cndmask_b32_e32 v10, v10, v18, vcc
	s_andn2_b64 s[0:1], s[0:1], exec
	s_branch BB1_8
BB1_23:                                 ; %.loopexit
	s_or_b64 exec, exec, s[6:7]
	v_mov_b32_e32 v18, 0
	v_lshlrev_b64 v[2:3], 4, v[17:18]
	v_add_u32_e32 v0, vcc, v0, v2
	v_addc_u32_e32 v1, vcc, v1, v3, vcc
	flat_store_dwordx4 v[0:1], v[5:8]
BB1_24:                                 ; %.loopexit7
	s_endpgm
.Lfunc_end1:
	.size	mergeSortPass, .Lfunc_end1-mergeSortPass
                                        ; -- End function
	.section	.AMDGPU.csdata
; Kernel info:
; codeLenInByte = 1680
; NumSgprs: 40
; NumVgprs: 37
; ScratchSize: 0
; MemoryBound: 0
; FloatMode: 192
; IeeeMode: 1
; LDSByteSize: 0 bytes/workgroup (compile time only)
; SGPRBlocks: 4
; VGPRBlocks: 9
; NumSGPRsForWavesPerEU: 40
; NumVGPRsForWavesPerEU: 37
; Occupancy: 6
; WaveLimiterHint : 1
; COMPUTE_PGM_RSRC2:USER_SGPR: 8
; COMPUTE_PGM_RSRC2:TRAP_HANDLER: 0
; COMPUTE_PGM_RSRC2:TGID_X_EN: 1
; COMPUTE_PGM_RSRC2:TGID_Y_EN: 0
; COMPUTE_PGM_RSRC2:TGID_Z_EN: 0
; COMPUTE_PGM_RSRC2:TIDIG_COMP_CNT: 0
	.text
	.protected	mergepack       ; -- Begin function mergepack
	.globl	mergepack
	.p2align	8
	.type	mergepack,@function
	.amdgpu_hsa_kernel mergepack
mergepack:                              ; @mergepack
mergepack$local:
	.amd_kernel_code_t
		amd_code_version_major = 1
		amd_code_version_minor = 2
		amd_machine_kind = 1
		amd_machine_version_major = 8
		amd_machine_version_minor = 0
		amd_machine_version_stepping = 3
		kernel_code_entry_byte_offset = 256
		kernel_code_prefetch_byte_size = 0
		granulated_workitem_vgpr_count = 1
		granulated_wavefront_sgpr_count = 2
		priority = 0
		float_mode = 192
		priv = 0
		enable_dx10_clamp = 1
		debug_mode = 0
		enable_ieee_mode = 1
		enable_wgp_mode = 0
		enable_mem_ordered = 0
		enable_fwd_progress = 0
		enable_sgpr_private_segment_wave_byte_offset = 0
		user_sgpr_count = 8
		enable_trap_handler = 0
		enable_sgpr_workgroup_id_x = 1
		enable_sgpr_workgroup_id_y = 1
		enable_sgpr_workgroup_id_z = 0
		enable_sgpr_workgroup_info = 0
		enable_vgpr_workitem_id = 0
		enable_exception_msb = 0
		granulated_lds_size = 0
		enable_exception = 0
		enable_sgpr_private_segment_buffer = 1
		enable_sgpr_dispatch_ptr = 1
		enable_sgpr_queue_ptr = 0
		enable_sgpr_kernarg_segment_ptr = 1
		enable_sgpr_dispatch_id = 0
		enable_sgpr_flat_scratch_init = 0
		enable_sgpr_private_segment_size = 0
		enable_sgpr_grid_workgroup_count_x = 0
		enable_sgpr_grid_workgroup_count_y = 0
		enable_sgpr_grid_workgroup_count_z = 0
		enable_wavefront_size32 = 0
		enable_ordered_append_gds = 0
		private_element_size = 1
		is_ptr64 = 1
		is_dynamic_callstack = 0
		is_debug_enabled = 0
		is_xnack_enabled = 0
		workitem_private_segment_byte_size = 0
		workgroup_group_segment_byte_size = 0
		gds_segment_byte_size = 0
		kernarg_segment_byte_size = 96
		workgroup_fbarrier_count = 0
		wavefront_sgpr_count = 18
		workitem_vgpr_count = 5
		reserved_vgpr_first = 0
		reserved_vgpr_count = 0
		reserved_sgpr_first = 0
		reserved_sgpr_count = 0
		debug_wavefront_private_segment_offset_sgpr = 0
		debug_private_segment_buffer_sgpr = 0
		kernarg_segment_alignment = 4
		group_segment_alignment = 4
		private_segment_alignment = 4
		wavefront_size = 6
		call_convention = -1
		runtime_loader_kernel_symbol = 0
	.end_amd_kernel_code_t
; %bb.0:
	s_load_dwordx4 s[12:15], s[6:7], 0x20
	s_load_dword s0, s[6:7], 0x28
	s_load_dword s1, s[4:5], 0x4
	s_mov_b32 s10, s9
	s_ashr_i32 s11, s9, 31
	s_waitcnt lgkmcnt(0)
	s_and_b32 s1, s1, 0xffff
	s_mul_i32 s8, s8, s1
	s_add_i32 s0, s0, s8
	v_add_u32_e32 v0, vcc, s0, v0
	s_lshl_b64 s[0:1], s[10:11], 2
	s_add_u32 s2, s12, s0
	s_addc_u32 s3, s13, s1
	s_add_u32 s4, 0, 0
	s_addc_u32 s5, s9, 1
	s_ashr_i64 s[4:5], s[4:5], 30
	s_add_u32 s4, s12, s4
	s_addc_u32 s5, s13, s5
	s_load_dword s2, s[2:3], 0x0
	s_load_dword s3, s[4:5], 0x0
	s_waitcnt lgkmcnt(0)
	v_add_u32_e32 v1, vcc, s2, v0
	v_cmp_gt_i32_e32 vcc, s3, v1
	s_and_saveexec_b64 s[2:3], vcc
	s_cbranch_execz BB2_2
; %bb.1:
	s_load_dwordx8 s[4:11], s[6:7], 0x0
	s_waitcnt lgkmcnt(0)
	s_add_u32 s2, s8, s0
	s_addc_u32 s3, s9, s1
	s_load_dword s2, s[2:3], 0x0
	v_mov_b32_e32 v4, s5
	s_waitcnt lgkmcnt(0)
	s_lshl_b32 s2, s2, 2
	s_add_u32 s0, s10, s0
	s_addc_u32 s1, s11, s1
	s_load_dword s0, s[0:1], 0x0
	s_waitcnt lgkmcnt(0)
	s_add_i32 s2, s2, s0
	v_add_u32_e32 v2, vcc, s2, v0
	v_ashrrev_i32_e32 v3, 31, v2
	v_lshlrev_b64 v[2:3], 2, v[2:3]
	v_add_u32_e32 v2, vcc, s4, v2
	v_addc_u32_e32 v3, vcc, v4, v3, vcc
	flat_load_dword v3, v[2:3]
	v_ashrrev_i32_e32 v2, 31, v1
	v_lshlrev_b64 v[0:1], 2, v[1:2]
	v_mov_b32_e32 v2, s7
	v_add_u32_e32 v0, vcc, s6, v0
	v_addc_u32_e32 v1, vcc, v2, v1, vcc
	s_waitcnt vmcnt(0) lgkmcnt(0)
	flat_store_dword v[0:1], v3
BB2_2:
	s_endpgm
.Lfunc_end2:
	.size	mergepack, .Lfunc_end2-mergepack
                                        ; -- End function
	.section	.AMDGPU.csdata
; Kernel info:
; codeLenInByte = 260
; NumSgprs: 18
; NumVgprs: 5
; ScratchSize: 0
; MemoryBound: 0
; FloatMode: 192
; IeeeMode: 1
; LDSByteSize: 0 bytes/workgroup (compile time only)
; SGPRBlocks: 2
; VGPRBlocks: 1
; NumSGPRsForWavesPerEU: 18
; NumVGPRsForWavesPerEU: 5
; Occupancy: 10
; WaveLimiterHint : 1
; COMPUTE_PGM_RSRC2:USER_SGPR: 8
; COMPUTE_PGM_RSRC2:TRAP_HANDLER: 0
; COMPUTE_PGM_RSRC2:TGID_X_EN: 1
; COMPUTE_PGM_RSRC2:TGID_Y_EN: 1
; COMPUTE_PGM_RSRC2:TGID_Z_EN: 0
; COMPUTE_PGM_RSRC2:TIDIG_COMP_CNT: 0
	.ident	"clang version 11.0.0 (/src/external/llvm-project/clang 0383ad1cfb0a8e05b0a020e8632400194628b243)"
	.section	".note.GNU-stack"
	.addrsig
	.amd_amdgpu_isa "amdgcn-amd-amdhsa--gfx803"
	.amd_amdgpu_hsa_metadata
---
Version:         [ 1, 0 ]
Kernels:
  - Name:            mergeSortFirst
    SymbolName:      'mergeSortFirst@kd'
    Language:        OpenCL C
    LanguageVersion: [ 1, 2 ]
    Args:
      - Name:            input
        TypeName:        'float4*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       F32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            result
        TypeName:        'float4*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       F32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            listsize
        TypeName:        int
        Size:            4
        Align:           4
        ValueKind:       ByValue
        ValueType:       I32
        AccQual:         Default
      - Size:            8
        Align:           8
        ValueKind:       HiddenGlobalOffsetX
        ValueType:       I64
      - Size:            8
        Align:           8
        ValueKind:       HiddenGlobalOffsetY
        ValueType:       I64
      - Size:            8
        Align:           8
        ValueKind:       HiddenGlobalOffsetZ
        ValueType:       I64
      - Size:            8
        Align:           8
        ValueKind:       HiddenNone
        ValueType:       I8
        AddrSpaceQual:   Global
      - Size:            8
        Align:           8
        ValueKind:       HiddenNone
        ValueType:       I8
        AddrSpaceQual:   Global
      - Size:            8
        Align:           8
        ValueKind:       HiddenNone
        ValueType:       I8
        AddrSpaceQual:   Global
      - Size:            8
        Align:           8
        ValueKind:       HiddenMultiGridSyncArg
        ValueType:       I8
        AddrSpaceQual:   Global
    CodeProps:
      KernargSegmentSize: 80
      GroupSegmentFixedSize: 0
      PrivateSegmentFixedSize: 0
      KernargSegmentAlign: 8
      WavefrontSize:   64
      NumSGPRs:        11
      NumVGPRs:        10
      MaxFlatWorkGroupSize: 256
  - Name:            mergeSortPass
    SymbolName:      'mergeSortPass@kd'
    Language:        OpenCL C
    LanguageVersion: [ 1, 2 ]
    Args:
      - Name:            input
        TypeName:        'float4*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       F32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            result
        TypeName:        'float4*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       F32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            nrElems
        TypeName:        int
        Size:            4
        Align:           4
        ValueKind:       ByValue
        ValueType:       I32
        AccQual:         Default
      - Name:            threadsPerDiv
        TypeName:        int
        Size:            4
        Align:           4
        ValueKind:       ByValue
        ValueType:       I32
        AccQual:         Default
      - Name:            constStartAddr
        TypeName:        'int*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       I32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Size:            8
        Align:           8
        ValueKind:       HiddenGlobalOffsetX
        ValueType:       I64
      - Size:            8
        Align:           8
        ValueKind:       HiddenGlobalOffsetY
        ValueType:       I64
      - Size:            8
        Align:           8
        ValueKind:       HiddenGlobalOffsetZ
        ValueType:       I64
      - Size:            8
        Align:           8
        ValueKind:       HiddenNone
        ValueType:       I8
        AddrSpaceQual:   Global
      - Size:            8
        Align:           8
        ValueKind:       HiddenNone
        ValueType:       I8
        AddrSpaceQual:   Global
      - Size:            8
        Align:           8
        ValueKind:       HiddenNone
        ValueType:       I8
        AddrSpaceQual:   Global
      - Size:            8
        Align:           8
        ValueKind:       HiddenMultiGridSyncArg
        ValueType:       I8
        AddrSpaceQual:   Global
    CodeProps:
      KernargSegmentSize: 88
      GroupSegmentFixedSize: 0
      PrivateSegmentFixedSize: 0
      KernargSegmentAlign: 8
      WavefrontSize:   64
      NumSGPRs:        40
      NumVGPRs:        37
      MaxFlatWorkGroupSize: 256
  - Name:            mergepack
    SymbolName:      'mergepack@kd'
    Language:        OpenCL C
    LanguageVersion: [ 1, 2 ]
    Args:
      - Name:            orig
        TypeName:        'float*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       F32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            result
        TypeName:        'float*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       F32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            constStartAddr
        TypeName:        'int*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       I32
        AddrSpaceQual:   Constant
        AccQual:         Default
        IsConst:         true
      - Name:            nullElems
        TypeName:        'int*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       I32
        AddrSpaceQual:   Constant
        AccQual:         Default
        IsConst:         true
      - Name:            finalStartAddr
        TypeName:        'int*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       I32
        AddrSpaceQual:   Constant
        AccQual:         Default
        IsConst:         true
      - Size:            8
        Align:           8
        ValueKind:       HiddenGlobalOffsetX
        ValueType:       I64
      - Size:            8
        Align:           8
        ValueKind:       HiddenGlobalOffsetY
        ValueType:       I64
      - Size:            8
        Align:           8
        ValueKind:       HiddenGlobalOffsetZ
        ValueType:       I64
      - Size:            8
        Align:           8
        ValueKind:       HiddenNone
        ValueType:       I8
        AddrSpaceQual:   Global
      - Size:            8
        Align:           8
        ValueKind:       HiddenNone
        ValueType:       I8
        AddrSpaceQual:   Global
      - Size:            8
        Align:           8
        ValueKind:       HiddenNone
        ValueType:       I8
        AddrSpaceQual:   Global
      - Size:            8
        Align:           8
        ValueKind:       HiddenMultiGridSyncArg
        ValueType:       I8
        AddrSpaceQual:   Global
    CodeProps:
      KernargSegmentSize: 96
      GroupSegmentFixedSize: 0
      PrivateSegmentFixedSize: 0
      KernargSegmentAlign: 8
      WavefrontSize:   64
      NumSGPRs:        18
      NumVGPRs:        5
      MaxFlatWorkGroupSize: 256
...

	.end_amd_amdgpu_hsa_metadata
