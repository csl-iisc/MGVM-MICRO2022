	.text
	.hsa_code_object_version 2,1
	.hsa_code_object_isa 8,0,3,"AMD","AMDGPU"
	.protected	FindKeyWithDigest_Kernel ; -- Begin function FindKeyWithDigest_Kernel
	.globl	FindKeyWithDigest_Kernel
	.p2align	8
	.type	FindKeyWithDigest_Kernel,@function
	.amdgpu_hsa_kernel FindKeyWithDigest_Kernel
FindKeyWithDigest_Kernel:               ; @FindKeyWithDigest_Kernel
FindKeyWithDigest_Kernel$local:
	.amd_kernel_code_t
		amd_code_version_major = 1
		amd_code_version_minor = 2
		amd_machine_kind = 1
		amd_machine_version_major = 8
		amd_machine_version_minor = 0
		amd_machine_version_stepping = 3
		kernel_code_entry_byte_offset = 256
		kernel_code_prefetch_byte_size = 0
		granulated_workitem_vgpr_count = 5
		granulated_wavefront_sgpr_count = 12
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
		kernarg_segment_byte_size = 112
		workgroup_fbarrier_count = 0
		wavefront_sgpr_count = 98
		workitem_vgpr_count = 23
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
	s_load_dwordx4 s[12:15], s[6:7], 0x10
	s_waitcnt lgkmcnt(0)
	s_cmp_lt_i32 s14, 1
	s_cbranch_scc1 BB0_39
; %bb.1:
	v_cvt_f32_u32_e32 v1, s14
	s_load_dword s0, s[4:5], 0x4
	s_mov_b32 s1, 0xffff
	s_load_dwordx4 s[80:83], s[6:7], 0x0
	s_load_dwordx4 s[84:87], s[6:7], 0x20
	s_load_dwordx4 s[88:91], s[6:7], 0x30
	v_rcp_iflag_f32_e32 v1, v1
	s_lshl_b32 s11, s13, 3
	s_waitcnt lgkmcnt(0)
	s_and_b32 s0, s0, s1
	s_mul_i32 s8, s8, s0
	v_mul_f32_e32 v1, 0x4f800000, v1
	v_cvt_u32_f32_e32 v1, v1
	v_add_u32_e32 v0, vcc, s8, v0
	v_mul_lo_u32 v0, v0, s14
	s_mov_b32 s4, 0x8000
	v_mul_lo_u32 v2, v1, s14
	v_mul_hi_u32 v3, v1, s14
	s_movk_i32 s5, 0x80
	s_add_i32 s8, s11, 0xa679438e
	v_sub_u32_e32 v4, vcc, 0, v2
	v_cmp_eq_u32_e64 s[0:1], 0, v3
	v_cndmask_b32_e64 v2, v2, v4, s[0:1]
	v_mul_hi_u32 v2, v2, v1
	s_add_i32 s9, s11, 0xc33707d6
	s_add_i32 s10, s11, 0xfde5380c
	s_add_i32 s11, s11, 0xab9423a7
	v_add_u32_e32 v3, vcc, v2, v1
	v_subrev_u32_e32 v1, vcc, v2, v1
	v_cndmask_b32_e64 v5, v1, v3, s[0:1]
	v_mul_hi_u32 v1, v5, v0
	s_mov_b32 s15, 0xefcdab89
	s_mov_b32 s16, 0x10325476
	s_mov_b32 s17, 0x98badcfe
	v_mul_lo_u32 v2, v1, s14
	v_add_u32_e32 v3, vcc, 1, v1
	v_add_u32_e32 v4, vcc, -1, v1
	s_mov_b32 s18, 0xf8fa0bcc
	v_sub_u32_e32 v6, vcc, v0, v2
	v_cmp_ge_u32_e32 vcc, v0, v2
	v_cmp_le_u32_e64 s[0:1], s14, v6
	s_and_b64 s[0:1], s[0:1], vcc
	v_cndmask_b32_e64 v1, v1, v3, s[0:1]
	v_cndmask_b32_e32 v7, v4, v1, vcc
	v_mul_hi_u32 v1, v5, v7
	s_mov_b32 s19, 0xbcdb4dd9
	s_mov_b32 s20, 0xb18b7a77
	s_mov_b32 s21, 0xe549bb38
	v_mul_lo_u32 v2, v1, s14
	v_add_u32_e32 v3, vcc, 1, v1
	v_add_u32_e32 v4, vcc, -1, v1
	s_mov_b32 s22, 0x4787c62a
	v_sub_u32_e32 v6, vcc, v7, v2
	v_cmp_ge_u32_e32 vcc, v7, v2
	v_cmp_le_u32_e64 s[0:1], s14, v6
	s_and_b64 s[0:1], s[0:1], vcc
	v_cndmask_b32_e64 v1, v1, v3, s[0:1]
	v_cndmask_b32_e32 v2, v4, v1, vcc
	v_mul_hi_u32 v1, v5, v2
	s_mov_b32 s23, 0xa8304613
	s_mov_b32 s24, 0xfd469501
	s_mov_b32 s25, 0x698098d8
	v_mul_lo_u32 v3, v1, s14
	v_add_u32_e32 v4, vcc, 1, v1
	v_add_u32_e32 v6, vcc, -1, v1
	s_mov_b32 s26, 0x8b44f7af
	v_sub_u32_e32 v8, vcc, v2, v3
	v_cmp_ge_u32_e32 vcc, v2, v3
	v_cmp_le_u32_e64 s[0:1], s14, v8
	s_and_b64 s[0:1], s[0:1], vcc
	v_cndmask_b32_e64 v1, v1, v4, s[0:1]
	v_cndmask_b32_e32 v3, v6, v1, vcc
	v_mul_hi_u32 v1, v5, v3
	s_mov_b32 s27, 0xffff5bb1
	s_mov_b32 s28, 0x895cd7be
	s_mov_b32 s29, 0x6b901122
	v_mul_lo_u32 v4, v1, s14
	v_add_u32_e32 v6, vcc, -1, v1
	s_mov_b32 s30, 0xfd987193
	s_mov_b32 s31, 0x49b40821
	v_sub_u32_e32 v8, vcc, v3, v4
	v_cmp_ge_u32_e32 vcc, v3, v4
	v_cmp_le_u32_e64 s[0:1], s14, v8
	v_add_u32_e64 v4, s[2:3], 1, v1
	s_and_b64 s[0:1], s[0:1], vcc
	v_cndmask_b32_e64 v1, v1, v4, s[0:1]
	v_cndmask_b32_e32 v4, v6, v1, vcc
	v_mul_hi_u32 v1, v5, v4
	v_mul_lo_u32 v6, v2, s14
	s_mov_b32 s33, 0xf61e2562
	s_mov_b32 s34, 0xc040b340
	v_mul_lo_u32 v8, v1, s14
	v_add_u32_e32 v9, vcc, -1, v1
	s_mov_b32 s35, 0x265e5a51
	s_mov_b32 s36, 0xe9b6c7aa
	v_sub_u32_e32 v10, vcc, v4, v8
	v_cmp_le_u32_e32 vcc, s14, v10
	v_cmp_ge_u32_e64 s[0:1], v4, v8
	v_add_u32_e64 v8, s[2:3], 1, v1
	s_and_b64 vcc, vcc, s[0:1]
	v_cndmask_b32_e32 v1, v1, v8, vcc
	v_cndmask_b32_e64 v8, v9, v1, s[0:1]
	v_mul_hi_u32 v9, v5, v8
	v_mul_lo_u32 v10, v3, s14
	v_sub_u32_e32 v1, vcc, v7, v6
	s_mov_b32 s37, 0xd62f105d
	v_mul_lo_u32 v6, v9, s14
	v_sub_u32_e32 v2, vcc, v2, v10
	v_mul_lo_u32 v10, v4, s14
	s_mov_b32 s38, 0x2441453
	v_sub_u32_e32 v11, vcc, v8, v6
	v_cmp_le_u32_e32 vcc, s14, v11
	v_add_u32_e64 v11, s[0:1], -1, v9
	v_cmp_ge_u32_e64 s[0:1], v8, v6
	v_add_u32_e64 v6, s[2:3], 1, v9
	s_and_b64 vcc, vcc, s[0:1]
	v_cndmask_b32_e32 v6, v9, v6, vcc
	v_cndmask_b32_e64 v6, v11, v6, s[0:1]
	v_mul_hi_u32 v9, v5, v6
	v_mul_lo_u32 v11, v8, s14
	v_sub_u32_e32 v3, vcc, v3, v10
	s_mov_b32 s39, 0xd8a1e681
	v_mul_lo_u32 v10, v9, s14
	v_sub_u32_e32 v4, vcc, v4, v11
	v_mul_lo_u32 v11, v6, s14
	s_mov_b32 s40, 0xe7d3fbc8
	v_sub_u32_e32 v12, vcc, v6, v10
	v_cmp_le_u32_e32 vcc, s14, v12
	v_add_u32_e64 v12, s[0:1], -1, v9
	v_cmp_ge_u32_e64 s[0:1], v6, v10
	v_add_u32_e64 v10, s[2:3], 1, v9
	s_and_b64 vcc, vcc, s[0:1]
	v_cndmask_b32_e32 v9, v9, v10, vcc
	v_cndmask_b32_e64 v9, v12, v9, s[0:1]
	v_mul_hi_u32 v10, v5, v9
	v_sub_u32_e32 v5, vcc, v8, v11
	v_mul_lo_u32 v12, v9, s14
	s_mov_b32 s41, 0x21e1cde6
	v_mul_lo_u32 v8, v10, s14
	s_mov_b32 s42, 0xf4d50d87
	v_sub_u32_e32 v6, vcc, v6, v12
	s_mov_b32 s43, 0x455a14ed
	v_sub_u32_e32 v10, vcc, v9, v8
	v_add_u32_e64 v11, s[0:1], s14, v10
	v_cmp_ge_u32_e64 s[0:1], v9, v8
	v_cmp_le_u32_e32 vcc, s14, v10
	v_mul_lo_u32 v9, v7, s14
	v_subrev_u32_e64 v8, s[2:3], s14, v10
	s_and_b64 vcc, vcc, s[0:1]
	v_cndmask_b32_e32 v7, v10, v8, vcc
	v_cndmask_b32_e64 v7, v11, v7, s[0:1]
	v_lshlrev_b32_e32 v10, 16, v2
	v_lshlrev_b32_e32 v11, 8, v1
	s_mov_b32 s0, 0xc06010c
	v_perm_b32 v10, v10, v11, s0
	v_lshlrev_b32_e32 v11, 8, v5
	v_sub_u32_e32 v8, vcc, v0, v9
	v_lshlrev_b32_e32 v9, 24, v3
	s_mov_b32 s0, 0x40c0c00
	v_or_b32_e32 v9, v10, v9
	v_lshlrev_b32_e32 v10, 16, v6
	v_and_b32_e32 v11, 0xff00, v11
	v_perm_b32 v12, v7, v4, s0
	v_and_b32_e32 v10, 0xff0000, v10
	v_or_b32_e32 v11, v12, v11
	v_or_b32_e32 v10, v11, v10
	s_brev_b32 s2, 1
	s_mov_b32 s3, 0x800000
	v_or_b32_e32 v11, s2, v10
	v_or_b32_e32 v12, s3, v10
	v_or_b32_e32 v13, s4, v10
	v_or_b32_e32 v14, s5, v10
	s_mov_b32 s44, 0xa9e3e905
	s_mov_b32 s45, 0xfcefa3f8
	s_mov_b32 s46, 0x676f02d9
	s_mov_b32 s47, 0x8d2a4c8a
	s_mov_b32 s48, 0xfffa3942
	s_mov_b32 s49, 0x8771f681
	s_mov_b32 s50, 0x6d9d6122
	s_mov_b32 s51, 0xa4beea44
	s_mov_b32 s52, 0x4bdecfa9
	s_mov_b32 s53, 0xf6bb4b60
	s_mov_b32 s54, 0xbebfbc70
	s_mov_b32 s55, 0x289b7ec6
	s_mov_b32 s56, 0xeaa127fa
	s_mov_b32 s57, 0xd4ef3085
	s_mov_b32 s58, 0x4881d05
	s_mov_b32 s59, 0xd9d4d039
	s_mov_b32 s60, 0xe6db99e5
	s_mov_b32 s61, 0x1fa27cf8
	s_mov_b32 s62, 0xc4ac5665
	s_mov_b32 s63, 0xf4292244
	s_mov_b32 s64, 0x432aff97
	s_mov_b32 s65, 0xfc93a039
	s_mov_b32 s66, 0x655b59c3
	s_mov_b32 s67, 0x8f0ccc92
	s_mov_b32 s68, 0xffeff47d
	s_mov_b32 s69, 0x85845dd1
	s_mov_b32 s70, 0x6fa87e4f
	s_mov_b32 s71, 0xfe2ce6e0
	s_mov_b32 s72, 0xa3014314
	s_mov_b32 s73, 0x4e0811a1
	s_mov_b32 s74, 0xf7537e82
	s_mov_b32 s75, 0xbd3af235
	s_mov_b32 s76, 0x2ad7d2bb
	s_mov_b64 s[6:7], 0
                                        ; implicit-def: $sgpr78_sgpr79
	s_branch BB0_4
BB0_2:                                  ;   in Loop: Header=BB0_4 Depth=1
	s_or_b64 exec, exec, s[0:1]
	s_add_i32 s14, s14, -1
	v_cmp_eq_u32_e64 s[0:1], s14, 0
	s_andn2_b64 s[78:79], s[78:79], exec
	s_and_b64 s[0:1], s[0:1], exec
	v_add_u16_e32 v8, 1, v8
	v_add_u32_e32 v0, vcc, 1, v0
	s_or_b64 s[78:79], s[78:79], s[0:1]
BB0_3:                                  ; %Flow63
                                        ;   in Loop: Header=BB0_4 Depth=1
	s_or_b64 exec, exec, s[90:91]
	s_and_b64 s[0:1], exec, s[78:79]
	s_or_b64 s[6:7], s[0:1], s[6:7]
	s_andn2_b64 exec, exec, s[6:7]
	s_cbranch_execz BB0_39
BB0_4:                                  ; =>This Inner Loop Header: Depth=1
	v_cmp_gt_i32_e32 vcc, s12, v0
	s_or_b64 s[78:79], s[78:79], exec
	s_and_saveexec_b64 s[90:91], vcc
	s_cbranch_execz BB0_3
; %bb.5:                                ;   in Loop: Header=BB0_4 Depth=1
	v_or_b32_sdwa v17, v9, v8 dst_sel:DWORD dst_unused:UNUSED_PAD src0_sel:DWORD src1_sel:BYTE_0
	s_cmp_lt_i32 s13, 4
	s_mov_b64 s[92:93], 0
	s_cbranch_scc1 BB0_11
; %bb.6:                                ; %NodeBlock44
                                        ;   in Loop: Header=BB0_4 Depth=1
	s_cmp_gt_i32 s13, 5
	s_cbranch_scc0 BB0_13
; %bb.7:                                ; %NodeBlock42
                                        ;   in Loop: Header=BB0_4 Depth=1
	s_cmp_gt_i32 s13, 6
	s_cbranch_scc0 BB0_14
; %bb.8:                                ; %LeafBlock40
                                        ;   in Loop: Header=BB0_4 Depth=1
	s_cmp_eq_u32 s13, 7
	s_mov_b64 s[0:1], -1
	s_cbranch_scc0 BB0_10
; %bb.9:                                ;   in Loop: Header=BB0_4 Depth=1
	s_mov_b64 s[0:1], 0
BB0_10:                                 ; %Flow54
                                        ;   in Loop: Header=BB0_4 Depth=1
	s_mov_b64 s[94:95], 0
	s_and_b64 vcc, exec, s[94:95]
	v_mov_b32_e32 v15, v11
	s_cbranch_vccnz BB0_15
	s_branch BB0_16
BB0_11:                                 ;   in Loop: Header=BB0_4 Depth=1
	s_mov_b64 s[0:1], 0
                                        ; implicit-def: $vgpr15
	s_cbranch_execnz BB0_22
BB0_12:                                 ;   in Loop: Header=BB0_4 Depth=1
	v_mov_b32_e32 v16, v17
	s_branch BB0_33
BB0_13:                                 ;   in Loop: Header=BB0_4 Depth=1
	s_mov_b64 s[0:1], 0
                                        ; implicit-def: $vgpr15
	s_cbranch_execnz BB0_17
	s_branch BB0_21
BB0_14:                                 ;   in Loop: Header=BB0_4 Depth=1
	s_mov_b64 s[0:1], 0
	v_mov_b32_e32 v15, v11
	s_cbranch_execz BB0_16
BB0_15:                                 ;   in Loop: Header=BB0_4 Depth=1
	v_mov_b32_e32 v15, v12
BB0_16:                                 ; %Flow55
                                        ;   in Loop: Header=BB0_4 Depth=1
	s_mov_b64 s[94:95], 0
	s_and_b64 vcc, exec, s[94:95]
	s_cbranch_vccz BB0_21
BB0_17:                                 ; %NodeBlock38
                                        ;   in Loop: Header=BB0_4 Depth=1
	s_cmp_gt_i32 s13, 4
	s_mov_b64 s[94:95], -1
	s_cbranch_scc0 BB0_19
; %bb.18:                               ;   in Loop: Header=BB0_4 Depth=1
	s_mov_b64 s[94:95], 0
BB0_19:                                 ; %Flow
                                        ;   in Loop: Header=BB0_4 Depth=1
	s_andn2_b64 vcc, exec, s[94:95]
	v_mov_b32_e32 v15, v13
	s_cbranch_vccnz BB0_21
; %bb.20:                               ;   in Loop: Header=BB0_4 Depth=1
	v_mov_b32_e32 v15, v14
BB0_21:                                 ; %Flow56
                                        ;   in Loop: Header=BB0_4 Depth=1
	s_mov_b64 s[94:95], 0
	s_and_b64 vcc, exec, s[94:95]
	s_cbranch_vccz BB0_12
BB0_22:                                 ; %NodeBlock36
                                        ;   in Loop: Header=BB0_4 Depth=1
	s_cmp_gt_i32 s13, 1
	s_mov_b64 s[92:93], -1
	v_mov_b32_e32 v16, v17
	s_cbranch_scc0 BB0_28
; %bb.23:                               ; %NodeBlock34
                                        ;   in Loop: Header=BB0_4 Depth=1
	s_cmp_gt_i32 s13, 2
	s_cbranch_scc0 BB0_25
; %bb.24:                               ;   in Loop: Header=BB0_4 Depth=1
	v_or_b32_e32 v16, s2, v17
	s_mov_b64 s[92:93], 0
	s_andn2_b64 vcc, exec, s[92:93]
	s_cbranch_vccz BB0_26
	s_branch BB0_27
BB0_25:                                 ;   in Loop: Header=BB0_4 Depth=1
	s_mov_b64 s[92:93], -1
                                        ; implicit-def: $vgpr16
	s_andn2_b64 vcc, exec, s[92:93]
	s_cbranch_vccnz BB0_27
BB0_26:                                 ;   in Loop: Header=BB0_4 Depth=1
	v_or_b32_e32 v16, s3, v17
BB0_27:                                 ; %Flow50
                                        ;   in Loop: Header=BB0_4 Depth=1
	s_mov_b64 s[92:93], 0
	v_mov_b32_e32 v15, v10
BB0_28:                                 ; %Flow58
                                        ;   in Loop: Header=BB0_4 Depth=1
	s_andn2_b64 vcc, exec, s[92:93]
	s_mov_b64 s[92:93], 0
	s_cbranch_vccnz BB0_33
; %bb.29:                               ; %NodeBlock
                                        ;   in Loop: Header=BB0_4 Depth=1
	s_cmp_gt_i32 s13, 0
	s_mov_b64 s[92:93], -1
	s_cbranch_scc0 BB0_31
; %bb.30:                               ;   in Loop: Header=BB0_4 Depth=1
	v_or_b32_e32 v16, s4, v17
	s_mov_b64 s[92:93], 0
	v_mov_b32_e32 v15, v10
BB0_31:                                 ; %Flow60
                                        ;   in Loop: Header=BB0_4 Depth=1
	s_andn2_b64 vcc, exec, s[92:93]
	s_mov_b64 s[92:93], 0
	s_cbranch_vccnz BB0_33
; %bb.32:                               ; %LeafBlock
                                        ;   in Loop: Header=BB0_4 Depth=1
	v_cmp_ne_u32_e64 s[0:1], s13, 0
	s_mov_b64 s[92:93], -1
BB0_33:                                 ; %Flow57
                                        ;   in Loop: Header=BB0_4 Depth=1
	s_and_b64 vcc, exec, s[0:1]
	s_cbranch_vccz BB0_35
; %bb.34:                               ; %NewDefault
                                        ;   in Loop: Header=BB0_4 Depth=1
	s_mov_b64 s[92:93], 0
	v_mov_b32_e32 v16, v17
	v_mov_b32_e32 v15, v10
BB0_35:                                 ; %Flow62
                                        ;   in Loop: Header=BB0_4 Depth=1
	s_andn2_b64 vcc, exec, s[92:93]
	s_cbranch_vccnz BB0_37
; %bb.36:                               ;   in Loop: Header=BB0_4 Depth=1
	v_or_b32_e32 v16, s5, v17
	v_mov_b32_e32 v15, v10
BB0_37:                                 ;   in Loop: Header=BB0_4 Depth=1
	v_add_u32_e32 v17, vcc, 0xd76aa477, v16
	v_alignbit_b32 v17, v17, v17, 25
	v_sub_u32_e32 v20, vcc, s16, v17
	v_add_u32_e32 v18, vcc, s15, v17
	v_and_b32_e32 v19, s15, v18
	v_and_b32_e32 v20, s17, v20
	v_or_b32_e32 v19, v19, v20
	v_add_u32_e32 v19, vcc, v15, v19
	v_add_u32_e32 v19, vcc, s18, v19
	v_alignbit_b32 v19, v19, v19, 20
	v_add_u32_e32 v19, vcc, v19, v18
	v_bfi_b32 v20, v19, v18, s15
	v_add_u32_e32 v20, vcc, s19, v20
	v_alignbit_b32 v20, v20, v20, 15
	v_add_u32_e32 v20, vcc, v20, v19
	v_bfi_b32 v18, v20, v19, v18
	v_add_u32_e32 v18, vcc, s20, v18
	v_alignbit_b32 v18, v18, v18, 10
	v_add_u32_e32 v18, vcc, v18, v20
	v_bfi_b32 v21, v18, v20, v19
	v_add_u32_e32 v17, vcc, v17, v21
	v_add_u32_e32 v17, vcc, s21, v17
	v_alignbit_b32 v17, v17, v17, 25
	v_add_u32_e32 v17, vcc, v17, v18
	v_bfi_b32 v21, v17, v18, v20
	v_add_u32_e32 v19, vcc, v19, v21
	v_add_u32_e32 v19, vcc, s22, v19
	v_alignbit_b32 v19, v19, v19, 20
	v_add_u32_e32 v19, vcc, v19, v17
	v_bfi_b32 v21, v19, v17, v18
	v_add_u32_e32 v20, vcc, v20, v21
	v_add_u32_e32 v20, vcc, s23, v20
	v_alignbit_b32 v20, v20, v20, 15
	v_add_u32_e32 v20, vcc, v20, v19
	v_bfi_b32 v21, v20, v19, v17
	v_add_u32_e32 v18, vcc, v18, v21
	v_add_u32_e32 v18, vcc, s24, v18
	v_alignbit_b32 v18, v18, v18, 10
	v_add_u32_e32 v18, vcc, v18, v20
	v_bfi_b32 v21, v18, v20, v19
	v_add_u32_e32 v17, vcc, v17, v21
	v_add_u32_e32 v17, vcc, s25, v17
	v_alignbit_b32 v17, v17, v17, 25
	v_add_u32_e32 v17, vcc, v17, v18
	v_bfi_b32 v21, v17, v18, v20
	v_add_u32_e32 v19, vcc, v19, v21
	v_add_u32_e32 v19, vcc, s26, v19
	v_alignbit_b32 v19, v19, v19, 20
	v_add_u32_e32 v19, vcc, v19, v17
	v_bfi_b32 v21, v19, v17, v18
	v_add_u32_e32 v20, vcc, v20, v21
	v_add_u32_e32 v20, vcc, s27, v20
	v_alignbit_b32 v20, v20, v20, 15
	v_add_u32_e32 v20, vcc, v20, v19
	v_bfi_b32 v21, v20, v19, v17
	v_add_u32_e32 v18, vcc, v18, v21
	v_add_u32_e32 v18, vcc, s28, v18
	v_alignbit_b32 v18, v18, v18, 10
	v_add_u32_e32 v18, vcc, v18, v20
	v_bfi_b32 v21, v18, v20, v19
	v_add_u32_e32 v17, vcc, v17, v21
	v_add_u32_e32 v17, vcc, s29, v17
	v_alignbit_b32 v17, v17, v17, 25
	v_add_u32_e32 v17, vcc, v17, v18
	v_bfi_b32 v21, v17, v18, v20
	v_add_u32_e32 v19, vcc, v19, v21
	v_add_u32_e32 v19, vcc, s30, v19
	v_alignbit_b32 v19, v19, v19, 20
	v_add_u32_e32 v19, vcc, v19, v17
	v_bfi_b32 v21, v19, v17, v18
	v_add_u32_e32 v20, vcc, s8, v20
	v_add_u32_e32 v20, vcc, v20, v21
	v_alignbit_b32 v20, v20, v20, 15
	v_add_u32_e32 v20, vcc, v20, v19
	v_bfi_b32 v21, v20, v19, v17
	v_add_u32_e32 v18, vcc, v18, v21
	v_add_u32_e32 v18, vcc, s31, v18
	v_alignbit_b32 v18, v18, v18, 10
	v_add_u32_e32 v18, vcc, v18, v20
	v_bfi_b32 v21, v19, v18, v20
	v_add_u32_e32 v17, vcc, v15, v17
	v_add_u32_e32 v17, vcc, v17, v21
	v_add_u32_e32 v17, vcc, s33, v17
	v_alignbit_b32 v17, v17, v17, 27
	v_add_u32_e32 v17, vcc, v17, v18
	v_bfi_b32 v21, v20, v17, v18
	v_add_u32_e32 v19, vcc, v19, v21
	v_add_u32_e32 v19, vcc, s34, v19
	v_alignbit_b32 v19, v19, v19, 23
	v_add_u32_e32 v19, vcc, v19, v17
	v_bfi_b32 v21, v18, v19, v17
	v_add_u32_e32 v20, vcc, v20, v21
	v_add_u32_e32 v20, vcc, s35, v20
	v_alignbit_b32 v20, v20, v20, 18
	v_add_u32_e32 v20, vcc, v20, v19
	v_bfi_b32 v21, v17, v20, v19
	v_add_u32_e32 v18, vcc, v16, v18
	v_add_u32_e32 v18, vcc, v18, v21
	v_add_u32_e32 v18, vcc, s36, v18
	v_alignbit_b32 v18, v18, v18, 12
	v_add_u32_e32 v18, vcc, v18, v20
	v_bfi_b32 v21, v19, v18, v20
	v_add_u32_e32 v17, vcc, v17, v21
	v_add_u32_e32 v17, vcc, s37, v17
	v_alignbit_b32 v17, v17, v17, 27
	v_add_u32_e32 v17, vcc, v17, v18
	v_bfi_b32 v21, v20, v17, v18
	v_add_u32_e32 v19, vcc, v19, v21
	v_add_u32_e32 v19, vcc, s38, v19
	v_alignbit_b32 v19, v19, v19, 23
	v_add_u32_e32 v19, vcc, v19, v17
	v_bfi_b32 v21, v18, v19, v17
	v_add_u32_e32 v20, vcc, v20, v21
	v_add_u32_e32 v20, vcc, s39, v20
	v_alignbit_b32 v20, v20, v20, 18
	v_add_u32_e32 v20, vcc, v20, v19
	v_bfi_b32 v21, v17, v20, v19
	v_add_u32_e32 v18, vcc, v18, v21
	v_add_u32_e32 v18, vcc, s40, v18
	v_alignbit_b32 v18, v18, v18, 12
	v_add_u32_e32 v18, vcc, v18, v20
	v_bfi_b32 v21, v19, v18, v20
	v_add_u32_e32 v17, vcc, v17, v21
	v_add_u32_e32 v17, vcc, s41, v17
	v_alignbit_b32 v17, v17, v17, 27
	v_add_u32_e32 v17, vcc, v17, v18
	v_bfi_b32 v21, v20, v17, v18
	v_add_u32_e32 v19, vcc, s9, v19
	v_add_u32_e32 v19, vcc, v19, v21
	v_alignbit_b32 v19, v19, v19, 23
	v_add_u32_e32 v19, vcc, v19, v17
	v_bfi_b32 v21, v18, v19, v17
	v_add_u32_e32 v20, vcc, v20, v21
	v_add_u32_e32 v20, vcc, s42, v20
	v_alignbit_b32 v20, v20, v20, 18
	v_add_u32_e32 v20, vcc, v20, v19
	v_bfi_b32 v21, v17, v20, v19
	v_add_u32_e32 v18, vcc, v18, v21
	v_add_u32_e32 v18, vcc, s43, v18
	v_alignbit_b32 v18, v18, v18, 12
	v_add_u32_e32 v18, vcc, v18, v20
	v_bfi_b32 v21, v19, v18, v20
	v_add_u32_e32 v17, vcc, v17, v21
	v_add_u32_e32 v17, vcc, s44, v17
	v_alignbit_b32 v17, v17, v17, 27
	v_add_u32_e32 v17, vcc, v17, v18
	v_bfi_b32 v21, v20, v17, v18
	v_add_u32_e32 v19, vcc, v19, v21
	v_add_u32_e32 v19, vcc, s45, v19
	v_alignbit_b32 v19, v19, v19, 23
	v_add_u32_e32 v19, vcc, v19, v17
	v_bfi_b32 v21, v18, v19, v17
	v_add_u32_e32 v20, vcc, v20, v21
	v_add_u32_e32 v20, vcc, s46, v20
	v_alignbit_b32 v20, v20, v20, 18
	v_add_u32_e32 v20, vcc, v20, v19
	v_bfi_b32 v21, v17, v20, v19
	v_add_u32_e32 v18, vcc, v18, v21
	v_add_u32_e32 v18, vcc, s47, v18
	v_alignbit_b32 v18, v18, v18, 12
	v_add_u32_e32 v18, vcc, v18, v20
	v_xor_b32_e32 v21, v18, v20
	v_xor_b32_e32 v22, v21, v19
	v_add_u32_e32 v17, vcc, v17, v22
	v_add_u32_e32 v17, vcc, s48, v17
	v_alignbit_b32 v17, v17, v17, 28
	v_add_u32_e32 v17, vcc, v17, v18
	v_xor_b32_e32 v21, v17, v21
	v_add_u32_e32 v19, vcc, v19, v21
	v_add_u32_e32 v19, vcc, s49, v19
	v_alignbit_b32 v19, v19, v19, 21
	v_add_u32_e32 v19, vcc, v19, v17
	v_xor_b32_e32 v21, v19, v17
	v_xor_b32_e32 v22, v21, v18
	v_add_u32_e32 v20, vcc, v20, v22
	v_add_u32_e32 v20, vcc, s50, v20
	v_alignbit_b32 v20, v20, v20, 16
	v_add_u32_e32 v20, vcc, v20, v19
	v_xor_b32_e32 v21, v21, v20
	v_add_u32_e32 v18, vcc, s10, v18
	v_add_u32_e32 v18, vcc, v18, v21
	v_alignbit_b32 v18, v18, v18, 9
	v_add_u32_e32 v18, vcc, v18, v20
	v_xor_b32_e32 v21, v20, v19
	v_xor_b32_e32 v21, v21, v18
	v_add_u32_e32 v17, vcc, v15, v17
	v_add_u32_e32 v17, vcc, v17, v21
	v_add_u32_e32 v17, vcc, s51, v17
	v_alignbit_b32 v17, v17, v17, 28
	v_add_u32_e32 v17, vcc, v17, v18
	v_xor_b32_e32 v21, v18, v20
	v_xor_b32_e32 v21, v21, v17
	v_add_u32_e32 v19, vcc, v19, v21
	v_add_u32_e32 v19, vcc, s52, v19
	v_alignbit_b32 v19, v19, v19, 21
	v_add_u32_e32 v19, vcc, v19, v17
	v_xor_b32_e32 v21, v17, v18
	v_xor_b32_e32 v21, v21, v19
	v_add_u32_e32 v20, vcc, v20, v21
	v_add_u32_e32 v20, vcc, s53, v20
	v_alignbit_b32 v20, v20, v20, 16
	v_add_u32_e32 v20, vcc, v20, v19
	v_xor_b32_e32 v21, v19, v17
	v_xor_b32_e32 v21, v21, v20
	v_add_u32_e32 v18, vcc, v18, v21
	v_add_u32_e32 v18, vcc, s54, v18
	v_alignbit_b32 v18, v18, v18, 9
	v_xor_b32_e32 v21, v20, v19
	v_add_u32_e32 v18, vcc, v18, v20
	v_xor_b32_e32 v21, v21, v18
	v_add_u32_e32 v17, vcc, v17, v21
	v_add_u32_e32 v17, vcc, s55, v17
	v_alignbit_b32 v17, v17, v17, 28
	v_add_u32_e32 v17, vcc, v17, v18
	v_xor_b32_e32 v21, v18, v20
	v_xor_b32_e32 v21, v21, v17
	v_add_u32_e32 v19, vcc, v16, v19
	v_add_u32_e32 v19, vcc, v19, v21
	v_add_u32_e32 v19, vcc, s56, v19
	v_alignbit_b32 v19, v19, v19, 21
	v_add_u32_e32 v19, vcc, v19, v17
	v_xor_b32_e32 v21, v17, v18
	v_xor_b32_e32 v21, v21, v19
	v_add_u32_e32 v20, vcc, v20, v21
	v_add_u32_e32 v20, vcc, s57, v20
	v_alignbit_b32 v20, v20, v20, 16
	v_add_u32_e32 v20, vcc, v20, v19
	v_xor_b32_e32 v21, v19, v17
	v_xor_b32_e32 v21, v21, v20
	v_add_u32_e32 v18, vcc, v18, v21
	v_add_u32_e32 v18, vcc, s58, v18
	v_alignbit_b32 v18, v18, v18, 9
	v_add_u32_e32 v18, vcc, v18, v20
	v_xor_b32_e32 v21, v20, v19
	v_xor_b32_e32 v21, v21, v18
	v_add_u32_e32 v17, vcc, v17, v21
	v_add_u32_e32 v17, vcc, s59, v17
	v_alignbit_b32 v17, v17, v17, 28
	v_add_u32_e32 v17, vcc, v17, v18
	v_xor_b32_e32 v21, v18, v20
	v_xor_b32_e32 v21, v21, v17
	v_add_u32_e32 v19, vcc, v19, v21
	v_add_u32_e32 v19, vcc, s60, v19
	v_alignbit_b32 v19, v19, v19, 21
	v_add_u32_e32 v19, vcc, v19, v17
	v_xor_b32_e32 v21, v17, v18
	v_xor_b32_e32 v21, v21, v19
	v_add_u32_e32 v20, vcc, v20, v21
	v_add_u32_e32 v20, vcc, s61, v20
	v_alignbit_b32 v20, v20, v20, 16
	v_add_u32_e32 v20, vcc, v20, v19
	v_xor_b32_e32 v21, v19, v17
	v_xor_b32_e32 v21, v21, v20
	v_add_u32_e32 v18, vcc, v18, v21
	v_add_u32_e32 v18, vcc, s62, v18
	v_alignbit_b32 v18, v18, v18, 9
	v_add_u32_e32 v18, vcc, v18, v20
	v_not_b32_e32 v21, v19
	v_or_b32_e32 v21, v18, v21
	v_add_u32_e32 v16, vcc, v16, v17
	v_xor_b32_e32 v21, v21, v20
	v_add_u32_e32 v16, vcc, v16, v21
	v_add_u32_e32 v16, vcc, s63, v16
	v_alignbit_b32 v16, v16, v16, 26
	v_not_b32_e32 v17, v20
	v_add_u32_e32 v16, vcc, v16, v18
	v_or_b32_e32 v17, v16, v17
	v_xor_b32_e32 v17, v17, v18
	v_add_u32_e32 v17, vcc, v19, v17
	v_add_u32_e32 v17, vcc, s64, v17
	v_alignbit_b32 v17, v17, v17, 22
	v_add_u32_e32 v17, vcc, v17, v16
	v_not_b32_e32 v19, v18
	v_or_b32_e32 v19, v17, v19
	v_xor_b32_e32 v19, v19, v16
	v_add_u32_e32 v20, vcc, s11, v20
	v_add_u32_e32 v19, vcc, v20, v19
	v_alignbit_b32 v19, v19, v19, 17
	v_add_u32_e32 v19, vcc, v19, v17
	v_not_b32_e32 v20, v16
	v_or_b32_e32 v20, v19, v20
	v_xor_b32_e32 v20, v20, v17
	v_add_u32_e32 v18, vcc, v18, v20
	v_add_u32_e32 v18, vcc, s65, v18
	v_alignbit_b32 v18, v18, v18, 11
	v_add_u32_e32 v18, vcc, v18, v19
	v_not_b32_e32 v20, v17
	v_or_b32_e32 v20, v18, v20
	v_xor_b32_e32 v20, v20, v19
	v_add_u32_e32 v16, vcc, v16, v20
	v_add_u32_e32 v16, vcc, s66, v16
	v_alignbit_b32 v16, v16, v16, 26
	v_add_u32_e32 v16, vcc, v16, v18
	v_not_b32_e32 v20, v19
	v_or_b32_e32 v20, v16, v20
	v_xor_b32_e32 v20, v20, v18
	v_add_u32_e32 v17, vcc, v17, v20
	v_add_u32_e32 v17, vcc, s67, v17
	v_alignbit_b32 v17, v17, v17, 22
	v_add_u32_e32 v17, vcc, v17, v16
	v_not_b32_e32 v20, v18
	v_or_b32_e32 v20, v17, v20
	v_xor_b32_e32 v20, v20, v16
	v_add_u32_e32 v19, vcc, v19, v20
	v_add_u32_e32 v19, vcc, s68, v19
	v_alignbit_b32 v19, v19, v19, 17
	v_add_u32_e32 v19, vcc, v19, v17
	v_not_b32_e32 v20, v16
	v_or_b32_e32 v20, v19, v20
	v_add_u32_e32 v15, vcc, v15, v18
	v_xor_b32_e32 v20, v20, v17
	v_add_u32_e32 v15, vcc, v15, v20
	v_add_u32_e32 v15, vcc, s69, v15
	v_alignbit_b32 v15, v15, v15, 11
	v_add_u32_e32 v15, vcc, v15, v19
	v_not_b32_e32 v18, v17
	v_or_b32_e32 v18, v15, v18
	v_xor_b32_e32 v18, v18, v19
	v_add_u32_e32 v16, vcc, v16, v18
	v_add_u32_e32 v16, vcc, s70, v16
	v_alignbit_b32 v16, v16, v16, 26
	v_add_u32_e32 v16, vcc, v16, v15
	v_not_b32_e32 v18, v19
	v_or_b32_e32 v18, v16, v18
	v_xor_b32_e32 v18, v18, v15
	v_add_u32_e32 v17, vcc, v17, v18
	v_add_u32_e32 v17, vcc, s71, v17
	v_alignbit_b32 v17, v17, v17, 22
	v_add_u32_e32 v17, vcc, v17, v16
	v_not_b32_e32 v18, v15
	v_or_b32_e32 v18, v17, v18
	v_xor_b32_e32 v18, v18, v16
	v_add_u32_e32 v18, vcc, v19, v18
	v_add_u32_e32 v18, vcc, s72, v18
	v_alignbit_b32 v18, v18, v18, 17
	v_add_u32_e32 v18, vcc, v18, v17
	v_not_b32_e32 v19, v16
	v_or_b32_e32 v19, v18, v19
	v_xor_b32_e32 v19, v19, v17
	v_add_u32_e32 v15, vcc, v15, v19
	v_add_u32_e32 v15, vcc, s73, v15
	v_alignbit_b32 v15, v15, v15, 11
	v_add_u32_e32 v15, vcc, v15, v18
	v_not_b32_e32 v19, v17
	v_or_b32_e32 v19, v15, v19
	v_xor_b32_e32 v19, v19, v18
	v_add_u32_e32 v16, vcc, v16, v19
	v_add_u32_e32 v16, vcc, s74, v16
	v_alignbit_b32 v16, v16, v16, 26
	v_add_u32_e32 v16, vcc, v16, v15
	v_not_b32_e32 v19, v18
	v_or_b32_e32 v19, v16, v19
	v_xor_b32_e32 v19, v19, v15
	v_add_u32_e32 v17, vcc, v17, v19
	v_add_u32_e32 v17, vcc, s75, v17
	v_alignbit_b32 v17, v17, v17, 22
	v_add_u32_e32 v17, vcc, v17, v16
	v_not_b32_e32 v19, v15
	v_or_b32_e32 v19, v17, v19
	v_xor_b32_e32 v19, v19, v16
	v_add_u32_e32 v18, vcc, v18, v19
	v_add_u32_e32 v18, vcc, s76, v18
	v_alignbit_b32 v18, v18, v18, 17
	v_not_b32_e32 v19, v16
	v_add_u32_e32 v18, vcc, v18, v17
	v_or_b32_e32 v19, v18, v19
	v_xor_b32_e32 v19, v19, v17
	v_add_u32_e32 v15, vcc, v15, v19
	v_add_u32_e32 v15, vcc, 0xeb86d391, v15
	v_alignbit_b32 v15, v15, v15, 11
	v_add_u32_e32 v15, vcc, v18, v15
	v_add_u32_e32 v15, vcc, s15, v15
	v_add_u32_e32 v16, vcc, 0x67452301, v16
	v_cmp_eq_u32_e32 vcc, s80, v16
	v_cmp_eq_u32_e64 s[0:1], s81, v15
	s_and_b64 s[0:1], vcc, s[0:1]
	v_add_u32_e32 v15, vcc, s17, v18
	v_cmp_eq_u32_e32 vcc, s82, v15
	s_and_b64 s[0:1], vcc, s[0:1]
	v_add_u32_e32 v15, vcc, s16, v17
	v_cmp_eq_u32_e32 vcc, s83, v15
	s_and_b64 s[92:93], vcc, s[0:1]
	s_and_saveexec_b64 s[0:1], s[92:93]
	s_cbranch_execz BB0_2
; %bb.38:                               ;   in Loop: Header=BB0_4 Depth=1
	v_mov_b32_e32 v15, s84
	v_mov_b32_e32 v16, s85
	flat_store_dword v[15:16], v0
	v_lshlrev_b16_e32 v15, 8, v5
	v_lshlrev_b16_e32 v16, 8, v7
	v_or_b32_sdwa v15, v4, v15 dst_sel:DWORD dst_unused:UNUSED_PAD src0_sel:BYTE_0 src1_sel:DWORD
	v_or_b32_sdwa v16, v6, v16 dst_sel:WORD_1 dst_unused:UNUSED_PAD src0_sel:BYTE_0 src1_sel:DWORD
	v_lshlrev_b16_e32 v17, 8, v3
	v_or_b32_sdwa v16, v15, v16 dst_sel:DWORD dst_unused:UNUSED_PAD src0_sel:WORD_0 src1_sel:DWORD
	v_lshlrev_b16_e32 v15, 8, v1
	v_or_b32_sdwa v15, v8, v15 dst_sel:DWORD dst_unused:UNUSED_PAD src0_sel:BYTE_0 src1_sel:DWORD
	v_or_b32_sdwa v17, v2, v17 dst_sel:WORD_1 dst_unused:UNUSED_PAD src0_sel:BYTE_0 src1_sel:DWORD
	v_or_b32_sdwa v15, v15, v17 dst_sel:DWORD dst_unused:UNUSED_PAD src0_sel:WORD_0 src1_sel:DWORD
	v_mov_b32_e32 v17, s86
	v_mov_b32_e32 v18, s87
	v_mov_b32_e32 v19, s88
	flat_store_dwordx2 v[17:18], v[15:16]
	v_mov_b32_e32 v15, s80
	v_mov_b32_e32 v16, s81
	v_mov_b32_e32 v17, s82
	v_mov_b32_e32 v18, s83
	v_mov_b32_e32 v20, s89
	flat_store_dwordx4 v[19:20], v[15:18]
	s_branch BB0_2
BB0_39:                                 ; %.loopexit
	s_endpgm
.Lfunc_end0:
	.size	FindKeyWithDigest_Kernel, .Lfunc_end0-FindKeyWithDigest_Kernel
                                        ; -- End function
	.section	.AMDGPU.csdata
; Kernel info:
; codeLenInByte = 3852
; NumSgprs: 98
; NumVgprs: 23
; ScratchSize: 0
; MemoryBound: 0
; FloatMode: 192
; IeeeMode: 1
; LDSByteSize: 0 bytes/workgroup (compile time only)
; SGPRBlocks: 12
; VGPRBlocks: 5
; NumSGPRsForWavesPerEU: 98
; NumVGPRsForWavesPerEU: 23
; Occupancy: 8
; WaveLimiterHint : 1
; COMPUTE_PGM_RSRC2:USER_SGPR: 8
; COMPUTE_PGM_RSRC2:TRAP_HANDLER: 0
; COMPUTE_PGM_RSRC2:TGID_X_EN: 1
; COMPUTE_PGM_RSRC2:TGID_Y_EN: 0
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
  - Name:            FindKeyWithDigest_Kernel
    SymbolName:      'FindKeyWithDigest_Kernel@kd'
    Language:        OpenCL C
    LanguageVersion: [ 1, 2 ]
    Args:
      - Name:            searchDigest0
        TypeName:        uint
        Size:            4
        Align:           4
        ValueKind:       ByValue
        ValueType:       U32
        AccQual:         Default
      - Name:            searchDigest1
        TypeName:        uint
        Size:            4
        Align:           4
        ValueKind:       ByValue
        ValueType:       U32
        AccQual:         Default
      - Name:            searchDigest2
        TypeName:        uint
        Size:            4
        Align:           4
        ValueKind:       ByValue
        ValueType:       U32
        AccQual:         Default
      - Name:            searchDigest3
        TypeName:        uint
        Size:            4
        Align:           4
        ValueKind:       ByValue
        ValueType:       U32
        AccQual:         Default
      - Name:            keyspace
        TypeName:        int
        Size:            4
        Align:           4
        ValueKind:       ByValue
        ValueType:       I32
        AccQual:         Default
      - Name:            byteLength
        TypeName:        int
        Size:            4
        Align:           4
        ValueKind:       ByValue
        ValueType:       I32
        AccQual:         Default
      - Name:            valsPerByte
        TypeName:        int
        Size:            4
        Align:           4
        ValueKind:       ByValue
        ValueType:       I32
        AccQual:         Default
      - Name:            foundIndex
        TypeName:        'int*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       I32
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            foundKey
        TypeName:        'uchar*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       U8
        AddrSpaceQual:   Global
        AccQual:         Default
      - Name:            foundDigest
        TypeName:        'uint*'
        Size:            8
        Align:           8
        ValueKind:       GlobalBuffer
        ValueType:       U32
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
      KernargSegmentSize: 112
      GroupSegmentFixedSize: 0
      PrivateSegmentFixedSize: 0
      KernargSegmentAlign: 8
      WavefrontSize:   64
      NumSGPRs:        98
      NumVGPRs:        23
      MaxFlatWorkGroupSize: 256
...

	.end_amd_amdgpu_hsa_metadata
